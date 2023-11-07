# jniDebug
A tool to make Android JNI Debugging in VSCode easier

## Installation
```
go install github.com/marcuswu/jnidebug@latest
```

Add the following lines to your .vscode/launch.json as the last two lines inside the list of configurations:
```
// #lldbclient-generated-begin
// #lldbclient-generated-end
```

### Usage
jniDebug -device emulator-0000 -lldb /path/to/lldb -package your.package -activity main.activity -debug /path/to/binary -vscode /path/to/launch.json [OPTIONS]

jniDebug requires adb to operate. Ensure it is installed and in your path.

Options:
-config  The name to use for the VSCode run configuration
-port    The port number to listen on
-wait    If present, pause app execution and wait for the debugger
-verbose If present, output verbose logging

example:
```
jniDebug -device emulator-5554 -lldb /Users/mwu/Library/Android/sdk/ndk/26.1.10909125/toolchains/llvm/prebuilt/darwin-x86_64/lib/clang/17.0.2/lib/linux/x86_64/lldb-server -package com.digitaltorque.structed -activity com.digitaltorque.structed.MainActivity -debug debug/jni/x86_64/libgojni.so -vscode ./.vscode/launch.json -config "Go Mobile Debugging" -wait
```

Once LLDB is running on the device, select and run the added launch configuration in VSCode.

This command is long, but most of these options will not change between executions so it is easy to write a short shell script for launching jnidebug.

I may add functionality later to find and use intelligent defaults for options such as the LLDB path.
