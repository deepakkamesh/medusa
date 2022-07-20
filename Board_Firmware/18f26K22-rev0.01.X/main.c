/*

 */

#include "mcc_generated_files/mcc.h"
#include "../lib/nrf24_lib.h"
#include "../lib/handler_lib.h"
#include "../lib/dht11_lib.h"



void main(void) {
    // Initialize the device
    SYSTEM_Initialize();

    // If using interrupts in PIC18 High/Low Priority Mode you need to enable the Global High and Low Interrupts
    // If using interrupts in PIC Mid-Range Compatibility Mode you need to enable the Global and Peripheral Interrupts
    // Use the following macros to:

    // Enable the Global Interrupts
    INTERRUPT_GlobalInterruptEnable();

    // Disable the Global Interrupts
    //INTERRUPT_GlobalInterruptDisable();

    // Enable the Peripheral Interrupts
    INTERRUPT_PeripheralInterruptEnable();

    // Disable the Peripheral Interrupts
    //INTERRUPT_PeripheralInterruptDisable();
    TMR1_SetInterruptHandler(TimerInterruptHandler);

    InitRadio();

    while (1) {
        HandlePacketLoop();
        NOP();
    }
}

void ProcessActionRequest(uint8_t actionID, uint8_t * data) {
    uint8_t tmpHumidity[] = {0, 0};

    switch (actionID) {
        case ACTION_STATUS_LED:
            LED_SetLow();
            if (data[0]) {
                LED_SetHigh();
            }
            break;
        case ACTION_RELOAD_CONFIG:
            ReloadConfig();
            break;
        case ACTION_GET_TEMP_HUMIDITY:
            GetMockTempHumidity(tmpHumidity);
            SendData(ACTION_GET_TEMP_HUMIDITY, tmpHumidity, 2);
            break;
        case ACTION_RESET_DEVICE:
            RESET();
            break;
        default:
            SendError(ERR_NOT_IMPL);
    }
}
