#ifndef KEYCHAIN_KEYCHAIN_H
#define KEYCHAIN_KEYCHAIN_H

#include <Security/Security.h>

OSStatus kcSet(const char *service, const char *account, const void *data, int length);
OSStatus kcGet(const char *service, const char *account, void **out, int *length);
OSStatus kcDelete(const char *service, const char *account);

#endif