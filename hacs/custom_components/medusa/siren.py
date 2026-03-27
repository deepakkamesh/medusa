"""Platform for siren integration."""
from homeassistant.components.siren import SirenEntity
from homeassistant.components.siren.const import SirenEntityFeature

from .const import DOMAIN, ACTION_BUZZER

async def async_setup_platform(hass, config, async_add_entities, discovery_info=None):
    """Set up the siren platform."""
    if discovery_info is None:
        return

    hub = hass.data[DOMAIN]["hub"]
    config_data = hass.data[DOMAIN]["config"]
    
    entities = []
    
    for board in config_data.get("Boards", []):
        room = board.get("Room", "unknown")
        name = board.get("Name", "unknown")
        actions = board.get("Actions", [])
        
        if ACTION_BUZZER in actions:
            entities.append(MedusaSiren(hub, room, name))

    async_add_entities(entities)


class MedusaSiren(SirenEntity):
    """Representation of a Medusa Siren."""

    def __init__(self, hub, room, board_name):
        """Initialize the siren."""
        self._hub = hub
        self._room = room
        self._board_name = board_name
        
        self._attr_name = f"{room} {board_name} buzzer"
        self._attr_unique_id = f"{room}_{board_name}_buzzer"
        self._attr_supported_features = SirenEntityFeature.TURN_ON | SirenEntityFeature.TURN_OFF
        self._is_on = False
        
        self._attr_device_info = {
            "identifiers": {(DOMAIN, f"{room}_{board_name}")},
            "name": f"{room}_{board_name}",
            "manufacturer": "Medusa",
            "suggested_area": room,
        }

    @property
    def is_on(self):
        """Return true if siren is on."""
        return self._is_on

    async def async_turn_on(self, **kwargs) -> None:
        """Turn the siren on."""
        # ACTION_BUZZER data payload: [on_byte, hi_duration, lo_duration]
        # Duration defaults to 100 ms (0x00, 0x64)
        data = bytearray([1, 0, 100])
        self._hub.send_action(self._room, self._board_name, ACTION_BUZZER, data)
        self._is_on = True
        self.async_write_ha_state()

    async def async_turn_off(self, **kwargs) -> None:
        """Turn the siren off."""
        # ACTION_BUZZER data payload: [0, 0, 0]
        data = bytearray([0, 0, 0])
        self._hub.send_action(self._room, self._board_name, ACTION_BUZZER, data)
        self._is_on = False
        self.async_write_ha_state()
