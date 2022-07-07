/*

 */

#include "mcc_generated_files/mcc.h"
#include "handler.h"
#include "../lib/nrf24_lib.h"



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
      /*  // Add your application code
        static char val = 0;
        // sprintf((char*) bufferTX, "Hello Arduino", val);
        nrf24_send_rf_data(bufferTX, sizeof (bufferTX));

        // Check for transmit and ack payload.
        while ((nrf24_read_register(NRF24_MEM_STATUSS) & 0x60) != 0x60) {
        }
        nrf24_read_rf_data(bufferRX, sizeof (bufferRX));
        LED_Toggle();
        sprintf((char*) bufferTX, "%s", bufferRX);
        val++;
        __delay_ms(500); // LED_Toggle();*/
        NOP();
    }
}
