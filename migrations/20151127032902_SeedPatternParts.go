package main

import (
	"database/sql"
	"github.com/Bnei-Baruch/mms-file-manager/models"
	"log"
	"github.com/Bnei-Baruch/mms-file-manager/config"
)

// Up is executed when this migration is applied
func Up_20151127032902(txn *sql.Tx) {
	db := config.NewDB()
	models.New(db)
	defer db.Close()

	pts := models.PatternParts{
		models.PatternPart{Key:"lecturer", Value:`rav|norav`, },
		models.PatternPart{Key:"lang", Value:`[[:lower:]]{3,4}`, },
		models.PatternPart{Key:"name", Value:`[a-z\-\d]+`, },
		models.PatternPart{Key:"content_type", Value:`[[:lower:]]+`, },
		models.PatternPart{Key:"line", Value:`[a-z\-\d]+`, },
		models.PatternPart{Key:"ot", Value:`o|t`, },
		models.PatternPart{Key:"date", Value:`\d{4}-\d{2}-\d{2}`, },
		models.PatternPart{Key:"cam", Value:`cam\d*_\d|xdcam\d*_\d{2,3}|cam\d*|xdcam\d*`, },
		models.PatternPart{Key:"archive_type", Value:`kabbalah|arvut|ligdol`, },
		models.PatternPart{Key:"index", Value:`n\d`, },
	}

	for _, pt := range pts {
		if err := pt.Save(); err != nil {
			log.Panicln("Unable to save pattern", pt, err)
		}
	}
}

// Down is executed when this migration is rolled back
func Down_20151127032902(txn *sql.Tx) {
	txn.Exec("delete from pattern_parts;")
}
