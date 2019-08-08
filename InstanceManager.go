package InstanceManager

import (
	"os"
	"Driscoll/system"
	"os/exec"
	"time"
	"strconv"
	"syscall"
	"strings"
	"Driscoll/io"
)

	func Run() {

		if(!system.ArgExists("instances") || system.GetArgAsInt("instances") < 1) {

			return
		}

		// Push ourselves into the background
		daemonise()

		// Kill any old sessions
		killOldSessions()

		time.Sleep(time.Second * 1)

		// Start the required number of new workers
		startWorkers()

		// Watch for an exit command
		go watchForExit()
	}

	func watchForExit() {

		// Watch for us being killed
		system.Exit()

		// Kill the workers
		killWorkers()

		// End ourselves
		os.Exit(0)
	}

	func startWorkers() {

		// Get our command without the instances argument
		var argsAsArray []string
		for key, value := range system.GetArgs() {

			if(key != "daemon" && key != "quiet" && key != "instances" && key != "manage") {

				argsAsArray = append(argsAsArray, "--" + key)
				argsAsArray = append(argsAsArray,value)
			}
		}

		// Enforce quiet mode on all workers
		argsAsArray = append(argsAsArray, "--quiet")

		// Run our command as many times as necessary
		for i := 0; i < system.GetArgAsInt("instances"); i++ {

			go monitorProcess(os.Args[0], argsAsArray)
		}
	}

	func monitorProcess(cmd string, args []string) {

		// Set up logging
		myLog := io.LogFile{}

		// Get the base application name
		bits := strings.Split(os.Args[0], "/")
		commandName := strings.Join(bits[len(bits) - 1:],"")

		// Set our log file
		myLog.Filename("log." + commandName)

		for {

			command := exec.Command(cmd,args...)
			crashBytes, _ := command.CombinedOutput()

			// Clean this up
			crash := strings.TrimSpace(string(crashBytes))

			if(len(crash) > 0) {

				myLog.Write("=== Service crashed. Output shown below. ===\n\n" + crash + "\n\n")
			}

			time.Sleep(time.Second * 1)
		}
	}

	func killOldSessions() {

		system.Exec("ps aux | grep " + os.Args[0] + " | grep -v grep | grep -v " + strconv.Itoa(os.Getpid()) + " | awk '{print $2}' | xargs kill")
	}

	func killWorkers() {

		system.Exec("ps aux | grep " + os.Args[0] + " | grep -v grep | grep -v " + strconv.Itoa(os.Getpid()) + " | grep -v instances | awk '{print $2}' | xargs kill")
	}

	func daemonise() {

		// If we are in the background, don't do this again, obviously
		if(system.ArgExists("manage")) {

			return
		}

		// Build u list of all arguments
		var argsAsArray []string

		// So we don't loop for ever, put a flag this is the manager running in the background
		argsAsArray = append(argsAsArray, "--manage")

		// Put the instances variable first
		argsAsArray = append(argsAsArray, "--instances")
		argsAsArray = append(argsAsArray, system.GetArg("instances"))

		for key, value := range system.GetArgs() {

			if(key != "daemon" && key != "quiet" && key != "instances") {

				argsAsArray = append(argsAsArray, "--" + key)
				argsAsArray = append(argsAsArray,value)
			}
		}

		cmd := exec.Command(os.Args[0],argsAsArray...)

		// Important: Make the child command a separate process group so it doesn't die with this program
		cmd.SysProcAttr = &syscall.SysProcAttr{ Setpgid:true }

		// Run the child process
		cmd.Start()

		// Kill this process
		os.Exit(0)
	}
