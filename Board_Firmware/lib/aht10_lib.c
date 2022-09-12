#include "aht10_lib.h"


bool AHT10Init(i2cAddr addr) {
    uint8_t cmd[3];
    uint8_t data[6];

    // Soft reset the device.
    cmd[0] = AHTX0_CMD_SOFTRESET;
    i2cWriteBytes(addr, cmd, 1);

    __delay_ms(20);

    // Send calibrate command. 
    cmd[0] = AHTX0_CMD_CALIBRATE;
    cmd[1] = 0x08;
    cmd[2] = 0x00;
    i2cWriteBytes(addr, cmd, 3);

    // Wait for ready.
    while (i2cRead1bReg(addr, 0x71) & AHTX0_STATUS_BUSY) {
        __delay_ms(10);
    }

    // Check if calibrated.
    if (!(i2cRead1bReg(addr, 0x71) & AHTX0_STATUS_CALIBRATED)) {
        return false;
    }
    return true;

}

void AHT10Read(i2cAddr addr, float *temp, float *humidity) {

    uint8_t cmd[3];
    uint8_t data[6];

    // Trigger reading.
    cmd[0] = AHTX0_CMD_TRIGGER;
    cmd[1] = 0x33;
    cmd[2] = 0x00;
    i2cWriteBytes(addr, cmd, 3);

    while (i2cRead1bReg(addr, 0x71) & AHTX0_STATUS_BUSY) {
        __delay_ms(10);
    }

    i2cReadBytes(addr, data, 6);

    uint32_t val = data[1];
    val <<= 8;
    val |= data[2];
    val <<= 4;
    val |= data[3] >> 4;
    *humidity = ((float) val * 100) / 0x100000;

    val  = data[3] & 0x0F;
    val  <<= 8;
    val  |= data[4];
    val  <<= 8;
    val  |= data[5];
    *temp = ((float) val  * 200 / 0x100000) - 50;
}
