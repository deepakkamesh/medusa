#include <xc.h>

uint8_t GetmockTemp(void) {
    return 35;
}

uint8_t GetmockHumidity(void) {
    return 40;
}

uint8_t GetMockTempHumidity(uint8_t *temp) {
    temp[0] = 35;
    temp[1] = 40;
    return 1;
}