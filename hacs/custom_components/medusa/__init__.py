"""The Medusa Hardware Controller integration."""
import asyncio
import json
import logging
import os
import voluptuous as vol

from homeassistant.core import HomeAssistant
from homeassistant.helpers.typing import ConfigType
import homeassistant.helpers.config_validation as cv
from homeassistant.helpers.event import async_track_time_interval
from homeassistant.const import EVENT_HOMEASSISTANT_STOP
from datetime import timedelta

from .const import (
    DOMAIN, CONF_FILE_PATH, CONF_PORT, DEFAULT_PORT,
    ACTION_TEMP, ACTION_LIGHT, ACTION_VOLT
)
from .hub import MedusaHub

_LOGGER = logging.getLogger(__name__)

PLATFORMS = ["sensor", "binary_sensor", "button", "siren"]

MEDUSA_SCHEMA = vol.Schema({
    vol.Required(CONF_FILE_PATH): cv.string,
    vol.Optional(CONF_PORT, default=DEFAULT_PORT): cv.port,
})

CONFIG_SCHEMA = vol.Schema({
    DOMAIN: MEDUSA_SCHEMA
}, extra=vol.ALLOW_EXTRA)


async def async_setup(hass: HomeAssistant, config: ConfigType) -> bool:
    """Set up the Medusa component."""
    if DOMAIN not in config:
        return True

    conf = config[DOMAIN]
    config_file_path = conf[CONF_FILE_PATH]
    port = conf[CONF_PORT]

    if not os.path.exists(config_file_path):
        _LOGGER.error("Medusa configuration file %s not found", config_file_path)
        return False

    with open(config_file_path, "r") as f:
        try:
            medusa_config = json.load(f)
        except json.JSONDecodeError as err:
            _LOGGER.error("Error parsing Medusa configuration file: %s", err)
            return False

    hub = MedusaHub(hass, medusa_config, port)
    
    hass.data.setdefault(DOMAIN, {})
    hass.data[DOMAIN]["hub"] = hub
    hass.data[DOMAIN]["config"] = medusa_config

    # Start the network hub
    await hub.start()

    async def handle_api_action(call):
        addr = call.data.get("addr")
        action_id = call.data.get("action_id")
        data = call.data.get("data", [])
        hub.api_action(addr, action_id, data)

    async def handle_relay_config_mode(call):
        hwaddr = call.data.get("hwaddr")
        on = call.data.get("on", False)
        hub.api_relay_config_mode(hwaddr, on)

    async def handle_board_config(call):
        addr = call.data.get("addr")
        paddr = call.data.get("paddr")
        hwaddr = call.data.get("hwaddr")
        naddr = call.data.get("naddr")
        hub.api_board_config(addr, paddr, hwaddr, naddr)

    hass.services.async_register(DOMAIN, "action", handle_api_action)
    hass.services.async_register(DOMAIN, "relay_config_mode", handle_relay_config_mode)
    hass.services.async_register(DOMAIN, "board_config", handle_board_config)

    async def async_stop_server(event):
        await hub.stop()

    hass.bus.async_listen_once(EVENT_HOMEASSISTANT_STOP, async_stop_server)

    async def async_poll_sensors(now):
        """Poll specific sensors that require active polling."""
        for board in medusa_config.get("Boards", []):
            room = board.get("Room", "unknown")
            name = board.get("Name", "unknown")
            actions = board.get("Actions", [])
            if ACTION_TEMP in actions:
                hub.send_action(room, name, ACTION_TEMP, b"")
            if ACTION_LIGHT in actions:
                hub.send_action(room, name, ACTION_LIGHT, b"")
            if ACTION_VOLT in actions:
                hub.send_action(room, name, ACTION_VOLT, b"")

    # The Go app polled based on PingInt, but we'll use a standard 60 seconds
    unsub_poll = async_track_time_interval(hass, async_poll_sensors, timedelta(seconds=60))
    
    async def async_clean_poll(event):
        unsub_poll()
    
    hass.bus.async_listen_once(EVENT_HOMEASSISTANT_STOP, async_clean_poll)

    # Load platforms
    for platform in PLATFORMS:
        hass.async_create_task(
            hass.helpers.discovery.async_load_platform(
                platform, DOMAIN, {}, config
            )
        )

    return True
