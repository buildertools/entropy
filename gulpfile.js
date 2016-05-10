var gulp  = require('gulp');
var child = require('child_process');
var server = null;
var it = null;

gulp.task('default', ['build', 'printHelp', 'watch', 'spawn']);

gulp.task('watch', function() {
  gulp.watch('./**/*.go', ['build','printHelp','spawn']);
  gulp.watch('./**/*.go', doFmt)
  gulp.watch('./**/*.go', doVet)
});

function doFmt(e) {
  var p = child.spawnSync('go', ['fmt', e.path]);
  dumpStreams("go fmt:\n", "go fmt (stderr):\n", p);
  return p;
}
function doVet(e) {
  var p = child.spawnSync('go', ['vet', e.path]);
  dumpStreams("go vet:\n", "go vet (stderr):\n", p);
  return p;
}

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
  server = child.spawn('entropy', ['--debug','manage', 'arg1', 'arg2', 'arg3']);
  server.stderr.on('data', function(data) {
    process.stdout.write(data.toString());
  });
  server.stdout.on('data', function(data) {
    process.stdout.write(data.toString());
  });
});

function forkIt(outPrefix, errPrefix, bin, args) {
  if (it)
    it.kill();
  it = child.spawn(bin, args);
  it.stderr.on('data', function(data) {
    process.stdout.write(errPrefix + data.toString());
  });
  it.stdout.on('data', function(data) {
    process.stdout.write(outPrefix + data.toString());
  });

}

function dumpStreams(outPrefix, errPrefix, p) {
  process.stdout.write(outPrefix + p.stdout.toString());
  process.stdout.write(errPrefix + p.stderr.toString());
}
