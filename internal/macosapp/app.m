#include <AppKit/AppKit.h>
#import <Cocoa/Cocoa.h>
#include <Foundation/Foundation.h>

#include <stdlib.h>
#include <sys/sysctl.h>
#include "app.h"

#if __has_include("_cgo_export.h")
#include "_cgo_export.h"
#endif

static NSString *const robloxBundleID = @"com.roblox.RobloxPlayer";

static NSArray<NSRunningApplication *> *gameInstances(void) {
    NSMutableArray *result = [NSMutableArray array];

    for (NSRunningApplication *app in [NSRunningApplication runningApplicationsWithBundleIdentifier:robloxBundleID]) {
        if (!app.terminated && !goIsTrayProcess(app.processIdentifier)) {
            [result addObject:app];
        }
    }

    return result;
}

@interface MRDelegate : NSObject <NSApplicationDelegate, NSMenuDelegate>
@property(strong) NSStatusItem *statusItem;
@property(strong) NSDateFormatter *timeFormatter;
@end

@implementation MRDelegate

- (void)applicationDidFinishLaunching:(NSNotification *)notification {
    self.timeFormatter = [[NSDateFormatter alloc] init];
    self.timeFormatter.timeStyle = NSDateFormatterShortStyle;

    NSMenu *menu = [[NSMenu alloc] init];
    menu.delegate = self;
    menu.autoenablesItems = NO;

    NSImage *icon = [NSImage imageNamed:@"MenuBarIcon"];
    if (icon == nil) {
        icon = [NSImage imageWithSystemSymbolName:@"square.stack.3d.up.fill" accessibilityDescription:@"twinblox"];
    }
    icon.template = YES;
    icon.size = NSMakeSize(18, 18);

    self.statusItem = [NSStatusBar.systemStatusBar statusItemWithLength:NSVariableStatusItemLength];
    self.statusItem.button.image = icon;
    self.statusItem.menu = menu;
}

- (void)menuNeedsUpdate:(NSMenu *)menu {
    [menu removeAllItems];

    [self addLabel:@"Launch any game through Roblox website" to:menu];
    [menu addItem:NSMenuItem.separatorItem];

    NSArray<NSRunningApplication *> *instances = gameInstances();
    if (instances.count == 0) {
       [self addLabel:@"No Instances running" to:menu];
    }

    for (NSRunningApplication *app in instances) {
        [self addInstance:app to:menu];
    }

    [menu addItem:NSMenuItem.separatorItem];
    NSMenuItem *closeAll = [self addAction:@selector(closeAll:) title:@"Close All Instances" to:menu];
    closeAll.enabled = instances.count > 0;

    [menu addItemWithTitle:@"Quit Twinblox" action:@selector(terminate:) keyEquivalent:@"q"];
}

- (void)addLabel:(NSString *)title to:(NSMenu *)menu {
    NSMenuItem *item = [menu addItemWithTitle:title action:nil keyEquivalent:@""];
    item.enabled = NO;
}

- (NSMenuItem *)addAction:(SEL)action title:(NSString *)title to:(NSMenu *)menu {
    NSMenuItem *item = [menu addItemWithTitle:title action:action keyEquivalent:@""];
    item.target = self;
    return item;
}

- (void)addInstance:(NSRunningApplication *)game to:(NSMenu *)menu {
    NSMenu *actions = [[NSMenu alloc] init];
    NSMenuItem *close = [self addAction:@selector(closeInstance:) title:@"Close" to:actions];
    close.representedObject = game;

    NSMenuItem *item = [menu addItemWithTitle:[self titleFor:game] action:nil keyEquivalent:@""];
    item.submenu = actions;
}

- (NSString *)titleFor:(NSRunningApplication *)game {
    char *label = goInstanceLabel(game.processIdentifier);
    NSString *who = label != NULL ? @(label) : @"unknown user";
    free(label);

    NSString *since = game.launchDate != nil ? [self.timeFormatter stringFromDate:game.launchDate] : @"?";
    return [NSString stringWithFormat:@"%@ • pid %d • %@", who, game.processIdentifier, since];
}

- (void)closeInstance:(NSMenuItem *)sender {
    [sender.representedObject forceTerminate];
}

- (void)closeAll:(id)sender {
    for (NSRunningApplication *app in gameInstances()) {
        [app forceTerminate];
    }
}

- (void)application:(NSApplication *)application openURLs:(NSArray<NSURL *> *)urls {
    for (NSURL *url in urls) {
        goOpenURL((char *)url.absoluteString.UTF8String);
    }
}

- (void)applicationWillTerminate:(NSNotification *)notification {
    goQuit();
}

@end

void RunApp(void) {
    @autoreleasepool {
        static MRDelegate *delegate;
        delegate = [[MRDelegate alloc] init];

        NSApplication.sharedApplication.activationPolicy = NSApplicationActivationPolicyAccessory;
        NSApplication.sharedApplication.delegate = delegate;

        [NSApp run];
    }
}