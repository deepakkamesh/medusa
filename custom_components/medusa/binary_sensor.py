"""Platform for binary sensor integration."""
from homeassistant.components.binary_sensor import BinarySensorEntity, BinarySensorDeviceClass
from homeassistant.core import callback
from homeassistant.helpers.dispatcher import async_dispatcher_connect

from .const import (
    DOMAIN, 
    SIGNAL_MEDUSA_UPDATE, 
    ACTION_MOTION, 
    ACTION_DOOR
)

async def async_setup_platform(hass, config, async_add_entities, discovery_info=None):
    """Set up the binary sensor platform."""
    if discovery_info is None:
        return

    hub = hass.data[DOMAIN]["hub"]
    config_data = hass.data[DOMAIN]["config"]
    
    entities = []
    
    for board in config_data.get("Boards", []):
        room = board.get("Room", "unknown")
        name = board.get("Name", "unknown")
        actions = board.get("Actions", [])
        
        if ACTION_MOTION in actions:
            entities.append(MedusaBinarySensor(hub, room, name, "motion", BinarySensorDeviceClass.MOTION))
            
        if ACTION_DOOR in actions:
            entities.append(MedusaBinarySensor(hub, room, name, "door", BinarySensorDeviceClass.DOOR))

    async_add_entities(entities)


class MedusaBinarySensor(BinarySensorEntity):
    """Representation of a Medusa Binary Sensor."""

    def __init__(self, hub, room, board_name, sensor_type, device_class):
        """Initialize the binary sensor."""
        self._hub = hub
        self._room = room
        self._board_name = board_name
        self._sensor_type = sensor_type
        
        self._attr_name = f"{room} {board_name} {sensor_type}"
        self._attr_unique_id = f"{room}_{board_name}_{sensor_type}"
        self._attr_device_class = device_class
        self._state = False
        
        self._attr_device_info = {
            "identifiers": {(DOMAIN, f"{room}_{board_name}")},
            "name": f"{room}_{board_name}",
            "manufacturer": "Medusa",
            "suggested_area": room,
        }

    @property
    def is_on(self):
        """Return true if the binary sensor is on."""
        return self._state

    async def async_added_to_hass(self):
        """Register callbacks."""
        self.async_on_remove(
            async_dispatcher_connect(
                self.hass,
                SIGNAL_MEDUSA_UPDATE.format(self._room, self._board_name),
                self._handle_update,
            )
        )

    @callback
    def _handle_update(self, data):
        """Handle updated data from the Medusa Hub."""
        if self._sensor_type in data:
            self._state = data[self._sensor_type]
            self.async_write_ha_state()
