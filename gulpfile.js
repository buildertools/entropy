var gulp  = require('gulp');
var gutil = require('gulp-util');
var go = require('gulp-go-tools');
var child = require('child_process');
var server = null;
var it = null;

var package = 'github.com/buildertools/entropy';
var output = 'bin/entropy';

var paths = {
  gofiles: './**/*.go',
  bin: '/go/bin/entropy'
};

var builds = {
  darwin64: { GOOS: "darwin", GOARCH: "amd64" },
  darwin64_2: { GOOS: "darwin", GOARCH: "amd64" },
  linux64: { GOOS: "linux", GOARCH: "amd64" }
};

gulp.task('default', ['build', 'printHelp', 'watch', 'spawn']);

gulp.task('watch', function() {
  gulp.watch('./**/*.go', ['fmt','build','buildall','printHelp','spawn']);
});

gulp.task('buildall', function(cb) {
  var e = process.env;
  for (var k in builds){
    if (builds.hasOwnProperty(k)) {
      e.GOOS = builds[k].GOOS;
      e.GOARCH = builds[k].GOARCH;
      child.exec('go build -o ' + output + '-' + k + ' ' + package, { env: e }, function(err, stdout, stderr) {
        if (err) cb(err);
        gutil.log(stdout, stderr);
      });
    }
  }
  cb();
});

function forkIt(outPrefix, errPrefix, bin, args, opts) {
  if (it)
    it.kill();
  it = child.spawn(bin, args, opts);
  if (it.error != null) 
    throw it.error;
  it.stderr.on('data', function(data) {
    process.stdout.write(errPrefix + data.toString());
  });
  it.stdout.on('data', function(data) {
    process.stdout.write(outPrefix + data.toString());
  });
}


gulp.task('fmt', function(cb) {
  var files = [];
  var that = this;
  child.exec('go fmt', function(err, stdout, stderr) {
    if (err) cb(err);
    gutil.log(stdout, stderr);
    cb();
  });
});


gulp.task('build', function() {
  var p = child.spawnSync('go', ['install']);
  dumpStreams("Target:\n", "Target (stderr):\n", p); 
  return p;
});

gulp.task('printHelp', function() {
  forkIt('helpOut', 'helpErr', 'entropy', [])
});

gulp.task('spawn', function() {
  if (server)
    server.kill();
  server = child.spawn('entropy', ['manage', "--image", "alpine",'tcp://swarm:3376']);
  server.stderr.on('data', function(data) {
    process.stdout.write(data.toString());
  });
  server.stdout.on('data', function(data) {
    process.stdout.write(data.toString());
  });
});

function dumpStreams(outPrefix, errPrefix, p) {
  process.stdout.write(outPrefix + p.stdout.toString());
  process.stdout.write(errPrefix + p.stderr.toString());
}

//function doFmt(e) {
//  var p = child.spawnSync('go', ['fmt', e.path]);
//  dumpStreams("go fmt:\n", "go fmt (stderr):\n", p);
//  return p;
//}
//function doVet(e) {
//  var p = child.spawnSync('go', ['vet', e.path]);
//  dumpStreams("go vet:\n", "go vet (stderr):\n", p);
//  return p;
//}


