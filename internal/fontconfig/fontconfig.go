package fontconfig

import (
	"log"
	"os"
	"path/filepath"
	"text/template"

	"github.com/dom1torii/even-better-font-manager/internal/config"
)

// text template instead of xml encoding because we want it to be readable
// and have comments if user decides to look inside the file.
const configTmpl = `<?xml version='1.0'?>
<!DOCTYPE fontconfig SYSTEM 'fonts.dtd'>
<fontconfig>
	<!-- MANAGED BY EVEN BETTER FONT MANAGER -->
	<!-- Choose reset in the tool or delete this file to stop using custom font. -->

	<!-- Both file name and font name work, but devs recommend file name. -->
	<!-- Font file should be located in /game/core/panorama/fonts/ebfm-custom/ -->
	<fontpattern>{{.FontFile}}</fontpattern>

	<!--
		We match fonts we want to replace inside <test>,
		replace it with our font inside first <edit> and select font size
		inside second <edit>. Only font name works, file name doesnt.
	-->
	<match target="pattern">
		<test name="family" compare="contains" qual="any">
			<!-- Main font CS2 uses is Stratum2 -->
			<string>Stratum2</string>
			<!-- Noto Sans is also used in CS2, though rarely. -->
			<string>notosans</string>
		</test>
		<edit name="family" mode="prepend" binding="strong">
			<string>{{.FontName}}</string>
		</edit>
		<edit name="pixelsize" mode="assign">
			<times>
				<name>pixelsize</name>
				<double>{{printf "%.1f" .Size}}</double>
			</times>
		</edit>
	</match>

	<!-- fonts-conf documentation: https://fontconfig.pages.freedesktop.org/fontconfig/fontconfig-user.html -->
</fontconfig>`

func Apply(cfg *config.Config, fontName string, fontPath string, fontSize float64) {
	fontFile := filepath.Base(fontPath)
	data := struct {
		FontFile string
		FontName string
		Size     float64
	}{
		FontFile: fontFile,
		FontName: fontName,
		Size:     fontSize,
	}

	confDir := filepath.Join(cfg.General.CS2Path, "game/core/panorama/fonts/conf.d")
	confPath := filepath.Join(confDir, "42-repl-global.conf")

	f, err := os.Create(confPath)
	if err != nil {
		log.Println("Failed to create config file: ", err)
		return
	}
	defer f.Close()

	tmpl, err := template.New("fontconfig").Parse(configTmpl)
	if err != nil {
		log.Println("Failed to parse template: ", err)
		return
	}

	err = tmpl.Execute(f, data)
	if err != nil {
		log.Println("Failed to write template to file: ", err)
	}
}

func Reset(cfg *config.Config) {
	confDir := filepath.Join(cfg.General.CS2Path, "game/core/panorama/fonts/conf.d")
	confPath := filepath.Join(confDir, "42-repl-global.conf")

	if err := os.RemoveAll(confPath); err != nil {
		log.Println("Failed to delete fonts config: ", err)
	}
}
