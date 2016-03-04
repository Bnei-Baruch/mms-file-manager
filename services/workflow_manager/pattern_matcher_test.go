package workflow_manager_test

import (
	"testing"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/Bnei-Baruch/mms-file-manager/models"
	wm "github.com/Bnei-Baruch/mms-file-manager/services/workflow_manager"
	"github.com/Bnei-Baruch/mms-file-manager/utils"
)

var patterns = []struct {
	expectedPriority int
	expectedRegex    string
	model            models.Pattern
}{
	{
		expectedPriority: 1,
		expectedRegex: `([[:lower:]]{3,4})_(arutz)_(\d{4}-\d{2}-\d{2})_([[:lower:]]+)_([a-z\-\d]+)_([a-z\-\d]+).(mpg)`,
		model: models.Pattern{
			Name: "lang_arutz_yyyy-mm-dd_type_line_name.mpg",
			Parts: models.Pairs{
				{Key: "lang", },
				{Key: "archive_type", Value: "arutz"},
				{Key: "date", },
				{Key: "content_type", },
				{Key: "line", },
				{Key: "name", },
			},
			Extension: "mpg",
		},
	},
	{
		expectedPriority: 1,
		expectedRegex: `([[:lower:]]{3,4})_(arutz)_(\d{4}-\d{2}-\d{2})_([[:lower:]]+)_([a-z\-\d]+).(mpg)`,
		model: models.Pattern{
			Name: "lang_arutz_yyyy-mm-dd_type_line.mpg",
			Parts: models.Pairs{
				{Key: "lang", },
				{Key: "archive_type", Value: "arutz"},
				{Key: "date", },
				{Key: "content_type", },
				{Key: "line", },
			},
			Extension: "mpg",
		},
	},
	{
		expectedPriority: 2,
		expectedRegex: `([[:lower:]]{3,4})_(arvut)_(\d{4}-\d{2}-\d{2})_(rawmaterial)_([a-z\-\d]+)_([a-z\-\d]+)_(cam\d*_\d|xdcam\d*_\d{2,3}|cam\d*|xdcam\d*).([[:alnum:]]{3,4})`,
		model: models.Pattern{
			Name: "lang_arvut_yyyy-mm-dd_rawmaterial_line_name_cam|xdcam(*).",
			Parts: models.Pairs{
				{Key: "lang", },
				{Key: "archive_type", Value: "arvut"},
				{Key: "date", },
				{Key: "content_type", Value: "rawmaterial"},
				{Key: "line", },
				{Key: "name", },
				{Key: "cam", },
			},
			Extension: "[[:alnum:]]{3,4}",
		},
	},
	{
		expectedPriority: 2,
		expectedRegex: `([[:lower:]]{3,4})_(arvut)_(\d{4}-\d{2}-\d{2})_(rawmaterial)_([a-z\-\d]+)_([a-z\-\d]+).(mpg)`,
		model: models.Pattern{
			Name: "lang_arvut_yyyy-mm-dd_rawmaterial_line_name.mpg",
			Parts: models.Pairs{
				{Key: "lang", },
				{Key: "archive_type", Value: "arvut"},
				{Key: "date", },
				{Key: "content_type", Value: "rawmaterial"},
				{Key: "line", },
				{Key: "name", },
			},
			Extension: "mpg",
		},
	},
	{
		expectedPriority: 1,
		expectedRegex: `([[:lower:]]{3,4})_(arvut)_(\d{4}-\d{2}-\d{2})_([[:lower:]]+)_([a-z\-\d]+)_([a-z\-\d]+).(mpg)`,
		model: models.Pattern{
			Name: "lang_arvut_yyyy-mm-dd_type_line_name.mpg",
			Parts: models.Pairs{
				{Key: "lang", },
				{Key: "archive_type", Value: "arvut"},
				{Key: "date", },
				{Key: "content_type"},
				{Key: "line", },
				{Key: "name", },
			},
			Extension: "mpg",
		},
	},
	{
		expectedPriority: 1,
		expectedRegex: `([[:lower:]]{3,4})_(arvut)_(\d{4}-\d{2}-\d{2})_([[:lower:]]+)_([a-z\-\d]+).(mpg)`,
		model: models.Pattern{
			Name: "lang_arvut_yyyy-mm-dd_type_line.mpg",
			Parts: models.Pairs{
				{Key: "lang", },
				{Key: "archive_type", Value: "arvut"},
				{Key: "date", },
				{Key: "content_type"},
				{Key: "line", },
			},
			Extension: "mpg",
		},
	},
	{
		expectedPriority: 1,
		expectedRegex: `([[:lower:]]{3,4})_(ligdol)_(\d{4}-\d{2}-\d{2})_([[:lower:]]+)_([a-z\-\d]+)_([a-z\-\d]+).(mpg)`,
		model: models.Pattern{
			Name: "lang_ligdol_yyyy-mm-dd_type_line_name.mpg",
			Parts: models.Pairs{
				{Key: "lang", },
				{Key: "archive_type", Value: "ligdol"},
				{Key: "date", },
				{Key: "content_type"},
				{Key: "line", },
				{Key: "name", },
			},
			Extension: "mpg",
		},
	},
	{
		expectedPriority: 1,
		expectedRegex: `([[:lower:]]{3,4})_(ligdol)_(\d{4}-\d{2}-\d{2})_([[:lower:]]+)_([a-z\-\d]+).(mpg)`,
		model: models.Pattern{
			Name: "lang_ligdol_yyyy-mm-dd_type_line.mpg",
			Parts: models.Pairs{
				{Key: "lang", },
				{Key: "archive_type", Value: "ligdol"},
				{Key: "date", },
				{Key: "content_type"},
				{Key: "line", },
			},
			Extension: "mpg",
		},
	},
	{
		expectedPriority: 1,
		expectedRegex: `([[:lower:]]{3,4})_(o|t)_(rav|norav)_(\d{4}-\d{2}-\d{2})_(lesson)_([a-z\-\d]+)_([a-z\-\d]+).(mpg)`,
		model: models.Pattern{
			Name: "lang_o|t_rav|norav_yyyy-mm-dd_lesson_line_name.mpg",
			Parts: models.Pairs{
				{Key: "lang", },
				{Key: "ot"},
				{Key: "lecturer", },
				{Key: "date", },
				{Key: "content_type", Value: "lesson"},
				{Key: "line", },
				{Key: "name", },
			},
			Extension: "mpg",
		},
	},
	{
		expectedPriority: 1,
		expectedRegex: `([[:lower:]]{3,4})_(o|t)_(rav|norav)_(\d{4}-\d{2}-\d{2})_(lesson)_([a-z\-\d]+).(mpg)`,
		model: models.Pattern{
			Name: "lang_o|t_rav|norav_yyyy-mm-dd_lesson_line.mpg",
			Parts: models.Pairs{
				{Key: "lang", },
				{Key: "ot"},
				{Key: "lecturer", },
				{Key: "date", },
				{Key: "content_type", Value: "lesson"},
				{Key: "line", },
			},
			Extension: "mpg",
		},
	},
	{
		expectedPriority: 1,
		expectedRegex: `([[:lower:]]{3,4})_(o|t)_(rav|norav)_(\d{4}-\d{2}-\d{2})_(rawmaterial)_([a-z\-\d]+)_([a-z\-\d]+)_(cam\d*_\d|xdcam\d*_\d{2,3}|cam\d*|xdcam\d*).([[:alnum:]]{3,4})`,
		model: models.Pattern{
			Name: "lang_o|t_rav|norav_yyyy-mm-dd_rawmaterial_line_name_cam|xdcam(*).*",
			Parts: models.Pairs{
				{Key: "lang", },
				{Key: "ot"},
				{Key: "lecturer", },
				{Key: "date", },
				{Key: "content_type", Value: "rawmaterial"},
				{Key: "line", },
				{Key: "name", },
				{Key: "cam", },
			},
			Extension: "[[:alnum:]]{3,4}",
		},
	},
	{
		expectedPriority: 1,
		expectedRegex: `([[:lower:]]{3,4})_(o|t)_(rav|norav)_(\d{4}-\d{2}-\d{2})_(rawmaterial)_([a-z\-\d]+)_([a-z\-\d]+).(mpg)`,
		model: models.Pattern{
			Name: "lang_o|t_rav|norav_yyyy-mm-dd_rawmaterial_line_name.mpg",
			Parts: models.Pairs{
				{Key: "lang", },
				{Key: "ot"},
				{Key: "lecturer", },
				{Key: "date", },
				{Key: "content_type", Value: "rawmaterial"},
				{Key: "line", },
				{Key: "name", },
			},
			Extension: "mpg",
		},
	},
	{
		expectedPriority: 0,
		expectedRegex: `([[:lower:]]{3,4})_(o|t)_(rav|norav)_(\d{4}-\d{2}-\d{2})_([[:lower:]]+)_([a-z\-\d]+)_([a-z\-\d]+).(mpg)`,
		model: models.Pattern{
			Name: "lang_o|t_rav|norav_yyyy-mm-dd_type_line_name.mpg",
			Parts: models.Pairs{
				{Key: "lang", },
				{Key: "ot"},
				{Key: "lecturer", },
				{Key: "date", },
				{Key: "content_type"},
				{Key: "line", },
				{Key: "name", },
			},
			Extension: "mpg",
		},
	},
	{
		expectedPriority: 0,
		expectedRegex: `([[:lower:]]{3,4})_(o|t)_(rav|norav)_(\d{4}-\d{2}-\d{2})_([[:lower:]]+)_([a-z\-\d]+).(mpg)`,
		model: models.Pattern{
			Name: "lang_o|t_rav|norav_yyyy-mm-dd_type_line.mpg",
			Parts: models.Pairs{
				{Key: "lang", },
				{Key: "ot"},
				{Key: "lecturer", },
				{Key: "date", },
				{Key: "content_type"},
				{Key: "line", },
			},
			Extension: "mpg",
		},
	},
	{
		expectedPriority: 3,
		expectedRegex: `(mlt)_(o)_(rav|norav)_([a-z\-\d]+)_(\d{4}-\d{2}-\d{2})_(lesson).(mp4)`,
		model: models.Pattern{
			Name: "mlt_o_rav|norav_name_yyyy-mm-dd_lesson.mp4",
			Parts: models.Pairs{
				{Key: "lang", Value: "mlt"},
				{Key: "ot", Value: "o"},
				{Key: "lecturer", },
				{Key: "name", },
				{Key: "date", },
				{Key: "content_type", Value: "lesson"},
			},
			Extension: "mp4",
		},
	},
	{
		expectedPriority: 3,
		expectedRegex: `(mlt)_(o)_(rav|norav)_([a-z\-\d]+)_(\d{4}-\d{2}-\d{2})_(lesson)_(n\d).(mp4)`,
		model: models.Pattern{
			Name: "mlt_o_rav|norav_name_yyyy-mm-dd_lesson_n0-9.mp4",
			Parts: models.Pairs{
				{Key: "lang", Value: "mlt"},
				{Key: "ot", Value: "o"},
				{Key: "lecturer", },
				{Key: "name", },
				{Key: "date", },
				{Key: "content_type", Value: "lesson"},
				{Key: "index"},
			},
			Extension: "mp4",
		},
	},
	{
		expectedPriority: 3,
		expectedRegex: `(mlt)_(o)_(rav|norav)_([a-z\-\d]+)_(\d{4}-\d{2}-\d{2})_(lesson).(mpg)`,
		model: models.Pattern{
			Name: "mlt_o_rav|norav_name_yyyy-mm-dd_lesson.mpg",
			Parts: models.Pairs{
				{Key: "lang", Value: "mlt"},
				{Key: "ot", Value: "o"},
				{Key: "lecturer", },
				{Key: "name", },
				{Key: "date", },
				{Key: "content_type", Value: "lesson"},
			},
			Extension: "mpg",
		},
	},
	{
		expectedPriority: 3,
		expectedRegex: `(mlt)_(o)_(rav|norav)_([a-z\-\d]+)_(\d{4}-\d{2}-\d{2})_(lesson)_(n\d).(mpg)`,
		model: models.Pattern{
			Name: "mlt_o_rav|norav_name_yyyy-mm-dd_lesson_n0-9.mpg",
			Parts: models.Pairs{
				{Key: "lang", Value: "mlt"},
				{Key: "ot", Value: "o"},
				{Key: "lecturer"},
				{Key: "name", },
				{Key: "date", },
				{Key: "content_type", Value: "lesson"},
				{Key: "index"},
			},
			Extension: "mpg",
		},
	},
}

var fileNames = []string{
	"mlt_o_rav_achana_2015-10-27_lesson.mpg",
	"heb_arutz_2012-12-16_film_crossroads.mpg",
	"heb_arutz_2013-04-20_promo_luah-shidurim_klali.mpg",
	"heb_arvut_2012-04-01_rawmaterial_wannabe_pilot_xdcam1_10.mp4",
	"heb_arvut_2012-09-24_program_hadashot-arvut.mpg",
	"heb_arvut_2013-05-02_rawmaterial_clip_itzik-safot.mpg",
	"heb_arvut_2015-04-22_promo_maagal-haim_bekarov.mpg",
	"heb_ligdol_2011-03-17_program_tofsim-olam_peter-vezeev.mpg",
	"heb_ligdol_2012-12-12_logo_ligdol-bekeif.mpg",
	//				"heb_o_rav_bs-tes-01_2003-12-05_lesson_n2.mp4",
	//				"heb_o_rav_rb-1988-10-dalet-midot_2003-03-30_lesson.mp4",
	"mlt_o_rav_2015-03-11_lesson_congress_n3.mpg",
	"mlt_o_rav_2015-10-04_lesson_full.mpg",
	"mlt_o_rav_2015-10-25_program_radio103fm.mpg",
	"mlt_o_rav_2015-11-15_program_lc_zima-2015.mpg",
	"rus_o_rav_2015-09-20_clip_lc_visshiy-razum.mpg",
	"rus_o_rav_2015-11-04_rawmaterial_taini-vechnoy-knigi_matot-7_xdcam_03.mp4",
	"rus_o_rav_2015-11-12_rawmaterial_sihot_baikot.mpg",
}


func TestPatternSpec(t *testing.T) {
	db = utils.SetupSpec()
	SetDefaultFailureMode(FailureHalts)
	Convey("Setup", t, func() {
		Convey("Subject: Pattern matching", func() {
			var file *models.File
				preparePatterns()

			Convey("When one pattern matched", func() {

				Convey("It should attach pattern to file", func() {
					for _, fileName := range fileNames {
						file = &models.File{
							FileName: fileName,
							TargetDir: "targetDir",
							EntryPoint: "label",
							SourcePath: "path",
						}

						file.CreateVersion()
						err := wm.AttachToPattern(file)
						So(err, ShouldBeNil)
						So(file.PatternId, ShouldNotBeNil)
						So(file.Status, ShouldEqual, models.HAS_PATTERN)
					}

				})
				Convey("It should parse fields from file name according to pattern", func() {
					for _, fileName := range fileNames {
						file = &models.File{
							FileName: fileName,
							TargetDir: "targetDir",
							EntryPoint: "label",
							SourcePath: "path",
						}

						file.CreateVersion()
						err := wm.AttachToPattern(file)
						So(err, ShouldBeNil)
						So(file.Attributes, ShouldNotBeNil)
						So(len(file.Attributes), ShouldEqual, len(file.Pattern.Parts))
						for _, p := range file.Pattern.Parts {
							value, ok := file.Attributes[p.Key]
							So(ok, ShouldBeTrue)
							So(value, ShouldNotBeBlank)
							//TODO: Check that values are correct - some should be included in list of values (like line, content_type)
						}
					}
				})
			})
			Convey("When no pattern is matched", func() {
				Convey("It should set file state to NO_PATTERN", func() {

					fileNames := []string{
						"heb_o_rav_bs-tes-01_2003-12-05_lesson_n2.mp4",
						"heb_o_rav_rb-1988-10-dalet-midot_2003-03-30_lesson.mp4",
					}

					for _, fileName := range fileNames {
						file = &models.File{
							FileName: fileName,
							TargetDir: "targetDir",
							EntryPoint: "label",
							SourcePath: "path",
						}

						file.CreateVersion()
						err := wm.AttachToPattern(file)
						So(err, ShouldBeNil)
						So(file.PatternId.Valid, ShouldBeFalse)
						So(file.Status, ShouldEqual, models.NO_PATTERN)
					}

				})
			})

			Convey("When more than one pattern is matched", func() {
				Convey("It should set file state to MANY_PATTERNS", func() {
					model := &models.Pattern{
						Name: "duplicate_mlt_o_rav|norav_name_yyyy-mm-dd_lesson_n0-9.mpg",
						Parts: models.Pairs{
							{Key: "lang", Value: "mlt"},
							{Key: "ot", Value: "o|t"},
							{Key: "lecturer"},
							{Key: "name", },
							{Key: "date", },
							{Key: "content_type", Value: "lesson"},
							{Key: "index"},
						},
						Extension: "mpg",
					}
					err := model.Save()
					So(err, ShouldBeNil)

					file = &models.File{
						FileName: "mlt_o_rav_rabash_2015-03-11_lesson_n3.mpg",
						TargetDir: "targetDir",
						EntryPoint: "label",
						SourcePath: "path",
					}

					file.CreateVersion()
					err = wm.AttachToPattern(file)
					So(err, ShouldBeNil)
					So(file.PatternId.Valid, ShouldBeFalse)
					So(file.Status, ShouldEqual, models.MANY_PATTERNS)
				})
			})

		})

		Convey("Describe Pattern saving", func() {
				db.Exec("DELETE FROM patterns;")

			Convey("It should reject unknown PatternPart keys", func() {
				model := models.Pattern{
					Name: "lang_arutz_yyyy-mm-dd_type_line_name.mpg",
					Parts: models.Pairs{
						{Key: "Unknown key", },
					},
					Extension: "mpg",
				}

				err := model.Save()
				So(err, ShouldNotBeNil)
			})

			Convey("It should reject patterns with duplicate names", func() {
				var err error
				pattern := patterns[0].model
				err = pattern.Save()
				So(err, ShouldBeNil)

				pattern = patterns[0].model
				err = pattern.Save()
				So(err, ShouldNotBeNil)

			})

			Convey("It should create pattern record, and calculate regex and priority", func() {
				for _, pat := range patterns {
					err := pat.model.Save()
					So(err, ShouldBeNil)

					p := &models.Pattern{Name: pat.model.Name}
					err = p.FindOne()
					So(err, ShouldBeNil)

					So(p.Regexp.Regx.String(), ShouldEqual, pat.expectedRegex)
					So(p.Priority, ShouldEqual, pat.expectedPriority)
				}
			})

		})
	})
}
func preparePatterns() {
	db.Exec("DELETE FROM patterns;")
	for _, pat := range patterns {
		pat.model.Save()
	}
}

