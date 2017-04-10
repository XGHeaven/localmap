def osarch(os, arch)
  desc "build executabel for #{os} with #{arch}"
  task "build-#{os}-#{arch}" do |t|
    sh "GOOS=#{os} GOARCH=#{arch} go build -v -o dist/localmap-#{os}-#{arch}"
  end
  task "build-#{os}-all": ["build-#{os}-#{arch}"]
end

directory "dist"

desc "build all platform in many os and different arch"
task "build-all": ["build-darwin-all", "build-windows-all", "build-linux-all"]

osarch :darwin, :amd64
osarch :linux, :amd64
osarch :linux, "386"
osarch :windows, "386"
osarch :windows, :amd64

desc "clean all dest file"
task :clean do
  sh "rm -rf dist"
end

task :build do |t|
  sh "go build -v -o dist/localmap"
end

task :default => [:build]
