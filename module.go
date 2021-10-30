package wasm

import "debug/dwarf"

// A Module represents a parsed WASM module.
type Module struct {
	// Sections contains the sections in the parsed file, in the order they
	// appear in the file. A valid  but empty file will have zero sections.
	//
	// The items in the slice will be a mix of the SectionXXX types.
	Sections []Section
}

// DWARF returns the DWARF debug information for this WebAssembly module.
func (m *Module) DWARF() (*dwarf.Data, error) {
	var abbrev, aranges, frame, info, line, pubnames, ranges, str []byte
	for _, section := range m.Sections {
		if section, ok := section.(*SectionCustom); ok {
			switch section.SectionName {
			case ".debug_abbrev":
				abbrev = section.Payload
			case ".debug_aranges":
				aranges = section.Payload
			case ".debug_frame":
				frame = section.Payload
			case ".debug_info":
				info = section.Payload
			case ".debug_line":
				line = section.Payload
			case ".debug_pubnames":
				pubnames = section.Payload
			case ".debug_ranges":
				ranges = section.Payload
			case ".debug_str":
				str = section.Payload
			}
		}
	}
	return dwarf.New(abbrev, aranges, frame, info, line, pubnames, ranges, str)
}
