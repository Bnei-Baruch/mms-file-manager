var gulp = require("gulp");
var coffee = require("gulp-coffee");
var include = require("gulp-include");
var rev = require('gulp-rev');
var gulpif = require('gulp-if');
var sass = require("gulp-sass");
var uglify = require("gulp-uglify");
var minify = require("gulp-minify-css");
var shell = require('gulp-shell');
var dotenv = require('dotenv');

var assets = [
  'assets/javascripts/application.coffee',
  'assets/stylesheets/application.scss'
]

// Development
// -----------------------------------------------------------

gulp.task("assets:compile", function() {
  gulp.src(assets)
    .pipe(include())
    .pipe(gulpif('*.coffee', coffee({bare: true})))
    .pipe(gulpif('*.scss', sass()))
    .pipe(gulp.dest("public/assets"))
    .pipe(rev())
    .pipe(gulp.dest("public/assets")) // file with digest
    .pipe(rev.manifest({path: 'manifest.json'}))
    .pipe(gulp.dest("public/assets")) // manifest.json
});

gulp.task("server", ["assets:compile"], function() {
  gulp.src('').pipe(shell('gin main.go'));
  gulp.watch(['assets/**/*.coffee', 'assets/**/*.scss'],{ interval: 500 }, ["assets:compile"]);
});

gulp.task("test", ["assets:compile"], function() {
  gulp.src('').pipe(shell('ginkgo -r'));
});

gulp.task("watch", ["assets:compile"], function() {
  gulp.src('').pipe(shell('ginkgo watch -failFast=true -noisyPendings=false -progress=true -randomizeAllSpecs=true -r=true -cover=true -v=true -trace=true'));
});

gulp.task("WEB-UI", ["assets:compile"], function() {
    gulp.src('').pipe(shell('goconvey -packages=1 -excludedDirs=models,routes,templates,migrations,public,assets,node_modules,vendor,tmp,config,tasks,test_files'));
});


// Production
// -----------------------------------------------------------

gulp.task("assets:precompile", function() {
  gulp.src(assets)
    .pipe(include())
    .pipe(gulpif('*.coffee', coffee({bare: true})))
    .pipe(gulpif('*.scss', sass()))
    .pipe(gulpif('*.css', minify()))
    .pipe(gulpif('*.js', uglify()))
    .pipe(gulp.dest("public/assets")) // file without digest
    .pipe(rev())
    .pipe(gulp.dest("public/assets")) // file with digest
    .pipe(rev.manifest({path: 'manifest.json'}))
    .pipe(gulp.dest("public/assets")) // manifest.json
});
