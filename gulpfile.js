var gulp  = require('gulp');
var child = require('child_process');
var server = null;

gulp.task('default', ['build', 'watch', 'spawn']);

gulp.task('watch', function() {
  gulp.watch('./**/*.go', ['build','spawn']);
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

gulp.task('spawn', function() {
  if (server)
    server.kill();
  server = child.spawn('entropy', ['--debug','version']);
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
