$version = "v0.1.0-beta2"
$importPath = "github.com/xgheaven/localmap"

def osarch(os, arch, suffix = nil)
  desc "build executable for #{os} with #{arch}"
  filename = "localmap-#{os}-#{arch}"
  filename += ".#{suffix}" if suffix
  task "build-#{os}-#{arch}" do |t|
    sh "GOOS=#{os} GOARCH=#{arch} go build -v -ldflags \"-X #{$importPath}/util.Version=#{$version} -X #{$importPath}/util.DateTime=#{Time.now.strftime("%Y-%m-%d.%H:%M")}\" -o dist/#{filename}"
  end

  desc "zip #{os}-#{arch} file"
  task "pack-#{os}-#{arch}": ["build-#{os}-#{arch}"] do |t|
    sh "cd dist && zip -r9 #{filename}.zip #{filename}"
  end

  task "build-#{os}-all": ["build-#{os}-#{arch}"]
  task "pack-#{os}-all": ["pack-#{os}-#{arch}"]
end

directory "dist"

desc "build all platform in many os and different arch"
task "build-all": ["build-darwin-all", "build-windows-all", "build-linux-all"]

desc "pack all platform in many os and different arch after build-all"
task "pack-all": ["pack-darwin-all", "pack-windows-all", "pack-linux-all"]

osarch :darwin, :amd64
osarch :linux, :amd64
osarch :linux, "386"
osarch :windows, "386", :exe
osarch :windows, :amd64, :exe

desc "clean all dest file"
task :clean do
  sh "rm -rf dist"
end

task :build do |t|
  sh "go build -v -ldflags \"-X #{$importPath}/util.Version=#{$version} -X #{$importPath}/util.DateTime=#{Time.now.strftime("%Y-%m-%d.%H:%M")}\" -o dist/localmap"
end

task :default => [:build]
