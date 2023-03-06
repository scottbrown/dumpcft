package main

var Profile string
var Verbose bool
var Regions string
var OutputDir string

func init() {
	rootCmd.PersistentFlags().StringVarP(&Profile, "profile", "p", DEFAULT_PROFILE, "The AWS profile to use")
	rootCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "Shows debug output")
	rootCmd.PersistentFlags().StringVarP(&Regions, "regions", "r", "", "One or more comma-delimited regions to dump")
	rootCmd.PersistentFlags().StringVarP(&OutputDir, "output-dir", "o", DEFAULT_OUTPUT_DIR, "The directory where templates are persisted to disk.")
}
