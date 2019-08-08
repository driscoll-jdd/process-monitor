# process-monitor
This is a simple library which allows any go application to run very safely in production. Taking inspiration from services such as Docker, this library allows you to run X instances of your application - if any instance crashes or is killed, another is spun up immediately. If there is any output from your application instance when it ends, this is logged - allowing you to diagnose crashes and work on future releases.

This also allows for simple upgrades - just overwrite the binary even as it is executing, and re-run it in monitor mode and it will replace the monitor instance plus all the child instances immediately with the updated binary.
