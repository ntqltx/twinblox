#import <Foundation/Foundation.h>
#include "keychain.h"

// builds a search query for keychain call
// this is made for easier saved password search that belongs to given account
static NSMutableDictionary *query(const char *service, const char *account) {
    return [@{
        (id)kSecClass: (id)kSecClassGenericPassword,
        (id)kSecAttrService: @(service),
        (id)kSecAttrAccount: @(account),
    } mutableCopy];
}

// saves data and replacing anything already saved for an account
OSStatus kcSet(const char *service, const char *account, const void *data, int length) {
    NSMutableDictionary *q = query(service, account);
    // add function fails if the item already exists
    // soo we're deleting it from the keychain before adding
    SecItemDelete((__bridge CFDictionaryRef)q);

    q[(id)kSecValueData] = [NSData dataWithBytes:data length:length];
    // only readable while the mac is unlocked
    q[(id)kSecAttrAccessible] = (id)kSecAttrAccessibleWhenUnlocked;
    return SecItemAdd((__bridge CFDictionaryRef)q, NULL);
}

// reads the saved data for an account
OSStatus kcGet(const char *service, const char *account, void **data, int *length) {
    NSMutableDictionary *q = query(service, account);
    q[(id)kSecReturnData] = @YES;

    CFTypeRef result = NULL;
    OSStatus status = SecItemCopyMatching((__bridge CFDictionaryRef)q, &result);

    if (status != errSecSuccess) {
        return status;
    }

    // copy into malloc memory for go to read it
    NSData *found = CFBridgingRelease(result);
    *length = (int)found.length;
    *data = malloc(found.length);

    memcpy(*data, found.bytes, found.length);
    return errSecSuccess;
}

// removes the saved data for an account
OSStatus kcDelete(const char *service, const char *account) {
    return SecItemDelete((__bridge CFDictionaryRef)query(service, account));
}