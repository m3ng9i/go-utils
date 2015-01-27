/*A simple log package.

With this package, you can create a logger with custom time format and message layout. You can output the log messages to an io.Writer object. If the io.Writer object is a file, you can choose to generate a new log file hourly, daily or monthly, and a datetime will be added to the old log file's filename.

This package do not support mail log, but you can define a function by yourself and use the function to do some extra log processing work, including mail log.
*/
package log
