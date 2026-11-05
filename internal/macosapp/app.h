#ifndef MACOSAPP_APP_H
#define MACOSAPP_APP_H

#include <stdbool.h>

void RunApp(void);

bool goIsTrayProcess(int pid);
char *goInstanceLabel(int pid);

void goOpenURL(char *url);
void goQuit(void);

#endif