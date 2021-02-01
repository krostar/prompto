package usecase

//
// func (pcf promptConfigFile) compiledFilename() string {
// 	return strings.TrimSuffix(pcf.String(), filepath.Ext(pcf.String())) + ".bin"
// }
//
// func (pcf promptConfigFile) compile(cfg *promptConfig) error {
// 	binaryName := pcf.compiledFilename()
//
// 	binaryFile, err := os.OpenFile(binaryName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
// 	if err != nil {
// 		return fmt.Errorf("unable to create binary file %q: %w", binaryName, err)
// 	}
// 	defer binaryFile.Close()
//
// 	if err := msgpack.NewEncoder(binaryFile).UseJSONTag(true).Encode(&cfg); err != nil {
// 		return fmt.Errorf("unable to marshal config file: %w", err)
// 	}
//
// 	if err := binaryFile.Sync(); err != nil {
// 		return fmt.Errorf("unable to sync file: %w", err)
// 	}
//
// 	return nil
// }
