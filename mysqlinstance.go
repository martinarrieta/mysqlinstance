// Package mysqlinstance is to create mysql instances from go.
package mysqlinstance

import (
    "os"
    "path"
    "fmt"
    "log"
    "time"
    "os/exec" 
)


import "github.com/robfig/config"


// Var debug is the variable used to display the debug messages
var debug bool = false

// Var Starttime is the number of seconds to wait until check if 
// the mysql srated or not. This will change in the future for a loop
var Starttime int = 15

// Type MySQLInstance is the main object
type MySQLInstance struct {
    Configfile string
    
    Innodbforcercovery bool
}

func Setdebug(deb bool) {
    debug = deb
}

// Funt Debugln is a funtion that check one constant "debug"
// and print the message if that constant is true
func Debugln( msg string ) { 

    if debug {
        log.Println("DEBUG: " + msg)
    }
    
}

// Func Fileexist return true if a given file name exist.
func Fileexist(file string) bool {
    
    if stat, err := os.Stat(file); !os.IsNotExist(err) && stat.Mode().IsRegular() {
        Debugln("[Fileexist] File " + file + " exists.")
        return true
    }
    Debugln("[Fileexist] File " + file + " does not exists." )
    
    return false
}

// Func Direxist return true if a given directory name exist.
func Direxist(file string) bool {
    if stat, err := os.Stat(file); !os.IsNotExist(err) && stat.Mode().IsDir() {
        Debugln("[Direxist] File " + file + " exists.")
        return true
    }
    Debugln("[Direxist] File " + file + " does not exists.")
    return false
}

// Func New creates and return a MySQLInstance object
func New( ) *MySQLInstance {
	m := new(MySQLInstance)
    
    return m
}

// Func Isrunning return true if the MySQLInstance is running.
// It use mysqladmin command with the option "status" to check it.
func (m MySQLInstance) Isrunning() bool { 
    
    command := m.Getbin("mysqladmin") 
    arg := "--defaults-file=" + m.Configfile + ""
    cmd := exec.Command( command, arg, "status")
    
    _, err := cmd.CombinedOutput()
    if err != nil {
        Debugln("Instance is not running.")
        return false
    } else {
        Debugln("Instance is running.")
        return true
    }
}

// Func Stop stops the MySQLInstance.
func (m MySQLInstance) Stop() bool { 
    
    
    command := m.Getbin("mysqladmin") 
    arg := "--defaults-file=" + m.Configfile + ""
    cmd := exec.Command( command, arg, "shutdown")
    
    output, err := cmd.CombinedOutput()
    if err != nil {
        Debugln(fmt.Sprint(err) + ": " + string(output))
        return false
    } else {
        log.Println("Instance stopped...")
        return true
    }
}


// Func Stop teruns true if the MySQLInstance is running and false if not.
func (m MySQLInstance) Status() bool { 
    
    if m.Isrunning() {
        Debugln("Instance is running.")
        return true
    } else {
        Debugln("Instance is NOT running.")
        return false
    }
}


// Func Start starts the MySQLInstance.
func (m MySQLInstance) Start() bool { 
    
    command := m.Getbin("mysqld") 
    
    arg := "--defaults-file=" + m.Configfile + ""
    
    if ! m.Isrunning() {
        
        log.Println("Starting the instance...")
        cmd := exec.Command(command, arg)
        
        if m.Innodbforcercovery {
            cmd.Args = append(cmd.Args, "--innodb-force-recovery=6")
        }
        
        err := cmd.Start()
        time.Sleep(15 * time.Second)
        
        if m.Isrunning() {
            log.Println("Instance started correctly.")
            return true
        } else {
            Debugln("Error, instance not started.")
            Debugln(fmt.Sprint(err))
            return false
        }
        
    } else {
        Debugln("Instance already running.")
        return true
    }
}

func (m MySQLInstance) Getconfigoption(section, option string) string {
    
    c, _ := config.ReadDefault(m.Configfile)
    
    ret, _ := c.String(section, option)
    
    return ret
    
}


func (m MySQLInstance) Getbin(command string) string {
    
    basedir := m.Getconfigoption("mysqld", "basedir")
    
    if basedir == "" {
        basedir = "/usr"
    }
    
    bin := ""
    
    switch command {
    case "mysqld": 
        
        bin = path.Join(basedir, "bin", "mysqld")
        
        if _, err := os.Stat(path.Join(basedir, "bin", "mysqld")); !os.IsNotExist(err) {
            return bin
        } else if _, err := os.Stat(path.Join(basedir, "sbin", "mysqld")); !os.IsNotExist(err) {
            return path.Join(basedir, "sbin", "mysqld")
        } else {
            log.Fatalln("We couldn't find the mysqld file")
        }
        
        
    case "mysqladmin": return path.Join(basedir, "bin", "mysqladmin") 
    case "mysql_install_db": 
        
        bin = path.Join(basedir, "bin", "mysql_install_db") 
        
        if _, err := os.Stat(path.Join(basedir, "bin", "mysql_install_db")); !os.IsNotExist(err) {
            return bin
        } else if _, err := os.Stat(path.Join(basedir, "scripts", "mysql_install_db")); !os.IsNotExist(err) {
            return path.Join(basedir, "scripts", "mysql_install_db")
        } else {
            log.Fatalln("We couldn't find the mysql_install_db file")
        }
    }
    
    return bin
}


// Func Initialize initializes the MySQLInstance.
func (m MySQLInstance) Initialize() bool { 
    
    
    if m.Isrunning(){
        log.Fatalln( "Instance is runnig, you have to stop it and clean the datadir directory to initialize it.")
    }
    
    command := m.Getbin("mysql_install_db")
    
    datadir := m.Getconfigoption("mysqld", "datadir")
    argbasedir := "--basedir=" + m.Getconfigoption("mysqld", "basedir")
    argdatadir := "--datadir=" + datadir
    arguser    := "--user=" + m.Getconfigoption("mysqld", "user")
    
    
    
    if ! Direxist(datadir) {
        log.Fatalln("The directory " + datadir + " does not exist.")
    }
    
    if Direxist(datadir + "/mysql") {
        log.Fatalln("The direcotry " + datadir + " is not empty")
    }
    
    cmd := exec.Command( command, argbasedir, argdatadir, arguser)
    
    output, err := cmd.CombinedOutput()
    if err != nil {
        log.Fatalln(fmt.Sprint(err) + ": " + string(output))
        return false
    } else {
        log.Println("Instance initialized correctly...")
        return true
    }
}
