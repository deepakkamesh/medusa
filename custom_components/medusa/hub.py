"""Hub for connecting to Medusa hardware."""
import asyncio
import logging
import struct
from collections import defaultdict
from homeassistant.core import HomeAssistant
from homeassistant.helpers.dispatcher import async_dispatcher_send

from .const import (
    DEFAULT_PORT,
    SIGNAL_MEDUSA_UPDATE,
    PKT_TYPE_PING,
    PKT_TYPE_DATA,
    PKT_TYPE_ACTION_REQ,
    PKT_TYPE_CONFIG,
    PKT_TYPE_FINDER,
    PKT_TYPE_RELAY_CFG_REQ,
    PKT_TYPE_RELAY_CFG_RESP,
    PKT_TYPE_RELAY_ERROR,
    PKT_TYPE_RELAY_BOARD_DATA,
    ACTION_MOTION,
    ACTION_TEMP,
    ACTION_LIGHT,
    ACTION_DOOR,
    ACTION_VOLT,
    ACTION_GAS,
    ACTION_PRESSURE,
    ACTION_ALTITUDE,
    ACTION_BUZZER,
    ACTION_RESET,
)

_LOGGER = logging.getLogger(__name__)

DEF_PIPE_ADDRESS = [0x68, 0x65, 0x6C, 0x6C, 0x6F]  # "hello"

class MedusaDatagramProtocol(asyncio.DatagramProtocol):
    """UDP Server to handle Relay Configuration Requests."""

    def __init__(self, hub):
        self._hub = hub

    def connection_made(self, transport):
        self.transport = transport
        _LOGGER.debug(f"Medusa UDP Server started on port {self._hub.port}")

    def datagram_received(self, data, addr):
        _LOGGER.debug(f"Received UDP datagram from {addr}: {data.hex()}")
        if not data:
            return

        if data[0] != PKT_TYPE_RELAY_CFG_REQ or len(data) < 7:
            _LOGGER.error(f"{addr} - Unknown packet type or bad len")
            return

        hwaddr = list(data[1:7])
        relay = self._hub.get_relay_by_hwaddr(hwaddr)
        
        if not relay:
            _LOGGER.error(f"{addr} - HWAddr not found: {hwaddr}")
            return

        # Prepare configuration response
        # Translated from makePktTypeRelayCfgResp in protocol.go
        resp = bytearray()
        resp.append(PKT_TYPE_RELAY_CFG_RESP)
        resp.extend(relay.get("PAddr0", [])) # 5 bytes
        resp.extend(relay.get("PAddr1", [])) # 5 bytes
        resp.append(relay.get("PAddr2", [0])[0])
        resp.append(relay.get("PAddr3", [0])[0])
        resp.append(relay.get("PAddr4", [0])[0])
        resp.append(relay.get("PAddr5", [0])[0])
        resp.append(relay.get("PAddr6", [0])[0])
        resp.append(relay.get("Channel", 0))
        resp.extend(relay.get("Addr", [])) # 3 bytes

        self.transport.sendto(resp, addr)
        _LOGGER.debug(f"Sent UDP cfg resp to {addr}: {resp.hex()}")


class MedusaProtocol(asyncio.Protocol):
    """TCP Server for Board Data Exchange."""

    def __init__(self, hub):
        self._hub = hub
        self._buffer = bytearray()
        self._transport = None
        self._peername = None

    def connection_made(self, transport):
        self._transport = transport
        self._peername = transport.get_extra_info("peername")
        _LOGGER.debug(f"New TCP connection from {self._peername}")
        # Identify relay by IP (naive mapping to first match in config)
        ip = self._peername[0]
        self._relay = self._hub.get_relay_by_ip(ip)

        if self._relay:
            # Store transport for writing actions back to boards attached to this relay
            hw_addr = tuple(self._relay["HWAddr"])
            self._hub.active_connections[hw_addr] = transport
        else:
            _LOGGER.warning(f"Connection from unknown relay IP: {ip}")

    def data_received(self, data):
        _LOGGER.debug(f"Data received from {self._peername}: {data.hex()}")
        self._buffer.extend(data)
        
        # Split buffer into packets based on length-prefixed protocol
        while len(self._buffer) > 0:
            pkt_len = self._buffer[0]
            if len(self._buffer) < pkt_len + 1:
                # Need more data
                break
            
            # Extract packet payload
            pkt = self._buffer[1:pkt_len + 1]
            # Advance buffer
            self._buffer = self._buffer[pkt_len + 1:]
            
            # Process packet
            try:
                self._hub.translate_packet(pkt)
            except Exception as ex:
                _LOGGER.error(f"Error parsing packet: {ex}", exc_info=True)

    def connection_lost(self, exc):
        _LOGGER.debug(f"TCP connection closed: {self._peername}")
        if self._relay:
            hw_addr = tuple(self._relay["HWAddr"])
            if hw_addr in self._hub.active_connections:
                del self._hub.active_connections[hw_addr]


class MedusaHub:
    """Manages Medusa Hardware connections and states."""

    def __init__(self, hass: HomeAssistant, config: dict, port: int):
        self._hass = hass
        self._config = config
        self.port = port
        self.active_connections = {}  # tuple(HWAddr) -> asyncio.Transport
        self.tcp_server = None
        self.udp_transport = None

    async def start(self):
        """Start TCP and UDP servers."""
        _LOGGER.info(f"Starting Medusa Hub on port {self.port}")
        loop = asyncio.get_running_loop()

        # UDP Server
        self.udp_transport, _ = await loop.create_datagram_endpoint(
            lambda: MedusaDatagramProtocol(self),
            local_addr=("0.0.0.0", self.port)
        )

        # TCP Server
        self.tcp_server = await loop.create_server(
            lambda: MedusaProtocol(self),
            "0.0.0.0", self.port
        )

    async def stop(self):
        """Stop servers."""
        if self.tcp_server:
            self.tcp_server.close()
            await self.tcp_server.wait_closed()
        if self.udp_transport:
            self.udp_transport.close()

    def get_relay_by_hwaddr(self, hwaddr):
        for r in self._config.get("Relays", []):
            if tuple(r.get("HWAddr", [])) == tuple(hwaddr):
                return r
        return None

    def get_relay_by_ip(self, ip):
        # Medusa's original config didn't track Relay IP statically in core.cfg.json
        # Relay IP logic in Go was naive: it just matched via UDP. 
        # But we'll just track whatever relay connected since we assume 1 relay per IP.
        # This function might just return the first Relay for simplicity unless further configured.
        # As per the go implementation, the IP was registered during UDP connection broadcast.
        # But for generic mapping without UDP state, return the first one as default if only 1 relay
        # Actually, let's just make sure active_connections stores what we need anyway.
        if len(self._config.get("Relays", [])) > 0:
            return self._config["Relays"][0]
        return None

    def get_relay_by_paddr(self, paddr):
        paddr_tuple = tuple(paddr)
        for r in self._config.get("Relays", []):
            for k in ["PAddr0", "PAddr1", "PAddr2", "PAddr3", "PAddr4", "PAddr5", "PAddr6"]:
                if tuple(r.get(k, [])) == paddr_tuple:
                    return r
            # Also check if it's matching PAddr[0] logic correctly.
            # Simplified matching.
        return None

    def get_board_by_addr(self, addr):
        addr_tup = tuple(addr)
        for b in self._config.get("Boards", []):
            if tuple(b.get("Addr", [])) == addr_tup:
                return b
        return None

    def get_boards_by_room(self, room):
        return [b for b in self._config.get("Boards", []) if b.get("Room") == room]

    def translate_packet(self, pkt):
        """Translate a physical payload to state updates."""
        if len(pkt) == 0:
            return
        
        pkt_type = pkt[0]
        if pkt_type == PKT_TYPE_RELAY_ERROR:
            _LOGGER.error(f"Relay Error Packet received: error_code={pkt[1]}")
            return
            
        if len(pkt) < 10:
            _LOGGER.warning("Packet too short")
            return
            
        if pkt_type != PKT_TYPE_RELAY_BOARD_DATA:
            _LOGGER.warning(f"Unknown packet type {pkt_type}")
            return

        paddr = list(pkt[1:6])
        board_pkt_type = pkt[6]
        addr = list(pkt[7:10])

        board = self.get_board_by_addr(addr)
        if not board:
            _LOGGER.warning(f"Unknown board addr {addr}")
            return

        room = board["Room"]
        name = board["Name"]
        state_updates = {}

        if board_pkt_type == PKT_TYPE_FINDER:
            _LOGGER.info(f"Board Find packet from {room}_{name}")
            return
            
        elif board_pkt_type == PKT_TYPE_PING:
            _LOGGER.debug(f"Ping from {room}_{name}")
            state_updates["ping"] = True
            
        elif board_pkt_type == PKT_TYPE_DATA:
            if len(pkt) < 12:
                _LOGGER.warning("Bad data packet length")
                return
            
            action = pkt[10]
            err_code = pkt[11]
            data = pkt[12:]
            
            if err_code != 0x00:
                _LOGGER.error(f"Error Code {err_code} from {room}_{name}")
                return

            if action == ACTION_TEMP:
                if len(data) >= 8:
                    raw_temp = struct.unpack('<f', data[0:4])[0]
                    t = raw_temp * 1.8 + 32  # Convert to fahrenheit locally
                    h = struct.unpack('<f', data[4:8])[0]
                    if -40 <= t <= 185 and 0 <= h <= 100:
                        state_updates["temp"] = round(t, 2)
                        state_updates["humidity"] = round(h, 2)
            
            elif action == ACTION_MOTION:
                state_updates["motion"] = bool(data[0] == 1)
                
            elif action == ACTION_DOOR:
                state_updates["door"] = bool(data[0] == 1)
                
            elif action == ACTION_VOLT:
                x = (data[1] << 8) | data[0]
                if x > 0:
                    state_updates["volt"] = round((2.048 * 1023) / float(x), 3)
                    
            elif action == ACTION_LIGHT:
                x = (data[1] << 8) | data[0]
                state_updates["light"] = round(3.3 * float(x) / 1023, 2)
                
            elif action == ACTION_GAS:
                raw_gas = struct.unpack('<I', data[0:4])[0]
                state_updates["gas"] = raw_gas / 1000.0  # Kohms
                
            elif action == ACTION_PRESSURE:
                state_updates["pressure"] = round(struct.unpack('<f', data[0:4])[0], 2)
                
            elif action == ACTION_ALTITUDE:
                state_updates["altitude"] = round(struct.unpack('<f', data[0:4])[0], 2)
                
            else:
                _LOGGER.warning(f"Unknown Action {action} from {room}_{name}")
                
        else:
            _LOGGER.warning(f"Unknown board packet type {board_pkt_type}")
            return

        if state_updates:
            # Emit signal via Home Assistant core dispatcher
            async_dispatcher_send(
                self._hass, 
                SIGNAL_MEDUSA_UPDATE.format(room, name), 
                state_updates
            )

    def send_action(self, room: str, name: str, action_id: int, data: bytes):
        """Send an action back down to Medusa hardware."""
        # Find the board
        board = None
        for b in self._config.get("Boards", []):
            if b.get("Room") == room and b.get("Name") == name:
                board = b
                break
                
        if not board:
            _LOGGER.error(f"Board not found for {room}_{name}")
            return
            
        addr = board.get("Addr", [])
        paddr_match = board.get("PAddr", [])
        
        # Find connection (In a multi-relay setup, you'd find by Relay's PAddr map)
        relay_transport = None
        for hw_addr, transport in self.active_connections.items():
            # For simplicity, if we have active connections, we send it there.
            relay_transport = transport
            break

        if not relay_transport:
            _LOGGER.error("No active connections to relays")
            return
            
        # PktTypeActionReq
        pkt = bytearray()
        pkt.append(PKT_TYPE_RELAY_BOARD_DATA)
        pkt.extend(paddr_match)
        pkt.append(PKT_TYPE_ACTION_REQ)
        pkt.extend(addr)
        pkt.append(action_id)
        if data:
            pkt.extend(data)
            
        # prepend length byte
        final_pkt = bytearray([len(pkt)]) + pkt
        try:
            relay_transport.write(final_pkt)
            _LOGGER.debug(f"Action {action_id} sent to {room}_{name}: {final_pkt.hex()}")
        except Exception as ex:
            _LOGGER.error(f"Failed to send action: {ex}")

    def api_action(self, addr: list[int], action_id: int, data: list[int]):
        """Equivalent to /api/action endpoint."""
        board = self.get_board_by_addr(addr)
        if not board:
            _LOGGER.error(f"api_action: board not found for addr {addr}")
            return
            
        paddr = board.get("PAddr", [])
        
        # In multi-relay, find by Relay's PAddr map. For simple implementation, use open connection.
        relay_transport = None
        for hw_addr, transport in self.active_connections.items():
            relay_transport = transport
            break

        if not relay_transport:
            _LOGGER.error("api_action: No active connections to relays")
            return
            
        pkt = bytearray()
        pkt.append(PKT_TYPE_RELAY_BOARD_DATA)
        pkt.extend(paddr)
        pkt.append(PKT_TYPE_ACTION_REQ)
        pkt.extend(addr)
        pkt.append(action_id)
        if data:
            pkt.extend(data)
            
        final_pkt = bytearray([len(pkt)]) + pkt
        try:
            relay_transport.write(final_pkt)
            _LOGGER.debug(f"api_action {action_id} sent: {final_pkt.hex()}")
        except Exception as ex:
            _LOGGER.error(f"Failed to send api_action: {ex}")

    def api_relay_config_mode(self, hwaddr: list[int], on: bool):
        """Equivalent to /api/relayconfigmode endpoint."""
        relay_transport = self.active_connections.get(tuple(hwaddr))
        if not relay_transport:
            _LOGGER.error(f"api_relay_config_mode: relay not found or not connected for hwaddr {hwaddr}")
            return
            
        relay = self.get_relay_by_hwaddr(hwaddr)
        if not relay:
            return

        resp = bytearray()
        resp.append(PKT_TYPE_RELAY_CFG_RESP)
        if on:
            resp.extend(DEF_PIPE_ADDRESS)
        else:
            resp.extend(relay.get("PAddr0", []))
            
        resp.extend(relay.get("PAddr1", []))
        resp.append(relay.get("PAddr2", [0])[0])
        resp.append(relay.get("PAddr3", [0])[0])
        resp.append(relay.get("PAddr4", [0])[0])
        resp.append(relay.get("PAddr5", [0])[0])
        resp.append(relay.get("PAddr6", [0])[0])
        resp.append(relay.get("Channel", 0))
        resp.extend(relay.get("Addr", []))

        final_pkt = bytearray([len(resp)]) + resp
        try:
            relay_transport.write(final_pkt)
            _LOGGER.debug(f"api_relay_config_mode sent to {hwaddr}: {final_pkt.hex()}")
        except Exception as ex:
            _LOGGER.error(f"Failed to send relay config Mode: {ex}")

    def api_board_config(self, addr: list[int], paddr: list[int], hwaddr: list[int], naddr: list[int]):
        """Equivalent to /api/boardconfig endpoint."""
        relay_transport = self.active_connections.get(tuple(hwaddr))
        if not relay_transport:
            _LOGGER.error(f"api_board_config: relay not found or not connected for hwaddr {hwaddr}")
            return

        board = self.get_board_by_addr(naddr)
        if not board:
            _LOGGER.error(f"api_board_config: board not found for default naddr {naddr}")
            return

        pkt = bytearray()
        pkt.append(PKT_TYPE_RELAY_BOARD_DATA)
        pkt.extend(paddr)
        pkt.append(PKT_TYPE_CONFIG)
        pkt.extend(addr)
        pkt.append(board.get("ARD", 0))
        pkt.append(board.get("PingInt", 0))
        pkt.extend(board.get("PAddr", []))
        pkt.extend(board.get("Addr", []))

        final_pkt = bytearray([len(pkt)]) + pkt
        try:
            relay_transport.write(final_pkt)
            _LOGGER.debug(f"api_board_config sent to {hwaddr}: {final_pkt.hex()}")
        except Exception as ex:
            _LOGGER.error(f"Failed to send board config: {ex}")
