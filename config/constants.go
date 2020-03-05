package config

// ConfigFilePath is the full local path to the the general config file of a component
const ConfigFilePath = ConfigFolder + ConfigFileName

// ConfigFileName is the name of the Iterum config files. This is what we look for
const ConfigFileName = "config.ivc"

// ConfigFolder is the name of the folder where we stroe all behind the scenes iterum work (like .git)
const ConfigFolder = ".iterum/"
