package workflow_manager_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/Bnei-Baruch/mms-file-manager/models"
	wm "github.com/Bnei-Baruch/mms-file-manager/services/workflow_manager"
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

var _ = Describe("Pattern matching", func() {
	var file *models.File
	BeforeEach(func() {
		preparePatterns()
	})

	Context("When one pattern matched", func() {

		It("should attach pattern to file", func() {
			for _, fileName := range fileNames {
				file = &models.File{
					FileName: fileName,
					TargetDir: "targetDir",
					Label: "label",
					SourcePath: "path",
				}

				err := wm.AttachToPattern(file)
				l.Println("Filename: ", fileName)
				Ω(err).ShouldNot(HaveOccurred())
				Ω(file.PatternId).ShouldNot(BeNil())
				Ω(file.Status).Should(Equal(models.HAS_PATTERN))
			}

		})
		PIt("should parse fields from file name according to pattern", func() {
			for _, fileName := range fileNames {
				file = &models.File{
					FileName: fileName,
					TargetDir: "targetDir",
					Label: "label",
					SourcePath: "path",
				}

				err := wm.AttachToPattern(file)
				Ω(err).ShouldNot(HaveOccurred())
				Ω(file.Attributes).ShouldNot(BeNil())
/*
				for _, p := range file.Pattern.Parts {
//					p.Key == file.Attributes
				}
*/

				//TODO:
				// 1. keys match
				// 2. each key has value
				// 3. values are correct
			}
		})
	})
	Context("When no pattern is matched", func() {
		It("should set file state to NO_PATTERN", func() {

			fileNames := []string{
				"heb_o_rav_bs-tes-01_2003-12-05_lesson_n2.mp4",
				"heb_o_rav_rb-1988-10-dalet-midot_2003-03-30_lesson.mp4",
			}

			for _, fileName := range fileNames {
				file = &models.File{
					FileName: fileName,
					TargetDir: "targetDir",
					Label: "label",
					SourcePath: "path",
				}

				err := wm.AttachToPattern(file)
				Ω(err).ShouldNot(HaveOccurred())
				Ω(file.PatternId.Valid).Should(BeFalse())
				Ω(file.Status).Should(Equal(models.NO_PATTERN))
			}

		})
	})


	Context("When more than one pattern is matched", func() {
		It("should set file state to MANY_PATTERNS", func() {
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
			Ω(err).ShouldNot(HaveOccurred())

			file = &models.File{
				FileName: "mlt_o_rav_rabash_2015-03-11_lesson_n3.mpg",
				TargetDir: "targetDir",
				Label: "label",
				SourcePath: "path",
			}

			err = wm.AttachToPattern(file)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(file.PatternId.Valid).Should(BeFalse())
			Ω(file.Status).Should(Equal(models.MANY_PATTERNS))
		})
	})

})

var _ = Describe("Pattern saving", func() {
	BeforeEach(func() {
		db.Exec("DELETE FROM patterns;")
	})

	It("must reject unknown PatternPart keys", func() {
		model := models.Pattern{
			Name: "lang_arutz_yyyy-mm-dd_type_line_name.mpg",
			Parts: models.Pairs{
				{Key: "Unknown key", },
			},
			Extension: "mpg",
		}

		err := model.Save()
		Ω(err).Should(HaveOccurred())
	})

	It("must reject patterns with duplicate names", func() {
		var err error
		pattern := patterns[0].model
		err = pattern.Save()
		Ω(err).ShouldNot(HaveOccurred())

		pattern = patterns[0].model
		err = pattern.Save()
		Ω(err).Should(HaveOccurred())
	})

	It("must create pattern record, and calculate regex and priority", func() {
		for _, pat := range patterns {
			err := pat.model.Save()
			Ω(err).ShouldNot(HaveOccurred())

			p := &models.Pattern{Name: pat.model.Name}
			err = p.FindOne()
			Ω(err).ShouldNot(HaveOccurred())

			Ω(p.Regexp.Regx.String()).Should(Equal(pat.expectedRegex))
			Ω(p.Priority).Should(Equal(pat.expectedPriority))
		}
	})

})

func preparePatterns() {
	db.Exec("DELETE FROM patterns;")
	for _, pat := range patterns {
		pat.model.Save()
	}
}
