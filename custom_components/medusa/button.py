"""Platform for button integration."""
from homeassistant.components.button import ButtonEntity, ButtonDeviceClass

from .const import DOMAIN, ACTION_RESET

async def async_setup_platform(hass, config, async_add_entities, discovery_info=None):
    """Set up the button platform."""
    if discovery_info is None:
        return

    hub = hass.data[DOMAIN]["hub"]
    config_data = hass.data[DOMAIN]["config"]
    
    entities = []
    
    for board in config_data.get("Boards", []):
        room = board.get("Room", "unknown")
        name = board.get("Name", "unknown")
        actions = board.get("Actions", [])
        
        if ACTION_RESET in actions:
            entities.append(MedusaButton(hub, room, name, "reset", ButtonDeviceClass.RESTART))

    async_add_entities(entities)


class MedusaButton(ButtonEntity):
    """Representation of a Medusa Button."""

    def __init__(self, hub, room, board_name, button_type, device_class):
        """Initialize the button."""
        self._hub = hub
        self._room = room
        self._board_name = board_name
        self._button_type = button_type
        
        self._attr_name = f"{room} {board_name} {button_type}"
        self._attr_unique_id = f"{room}_{board_name}_{button_type}"
        self._attr_device_class = device_class
        
        self._attr_device_info = {
            "identifiers": {(DOMAIN, f"{room}_{board_name}")},
            "name": f"{room}_{board_name}",
            "manufacturer": "Medusa",
            "suggested_area": room,
        }

    async def async_press(self) -> None:
        """Handle the button press."""
        if self._button_type == "reset":
            # Medusa Reset Action: data is empty
            self._hub.send_action(self._room, self._board_name, ACTION_RESET, b"")
