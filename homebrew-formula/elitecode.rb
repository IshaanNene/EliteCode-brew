# homebrew-formula/elitecode.rb
class Elitecode < Formula
  desc "CLI tool for competitive programming and coding challenges"
  homepage "https://github.com/yourusername/elitecode"
  url "https://github.com/yourusername/elitecode/archive/v1.0.0.tar.gz"
  sha256 "YOUR_SHA256_HERE"
  license "MIT"

  depends_on "node" => :build
  depends_on "go" => :build
  depends_on "docker"

  def install
    # Build Go version
    cd "golang" do
      system "go", "build", "-o", "#{bin}/elitecode", "."
    end

    # Install Node.js version as fallback
    cd "nodejs" do
      system "npm", "install"
      system "npm", "run", "build"
      libexec.install Dir["*"]
      (bin/"elitecode-node").write_env_script("#{Formula["node"].bin}/node", "#{libexec}/dist/index.js", {})
    end

    # Install templates and configs
    (etc/"elitecode").install "templates"
    (etc/"elitecode").install "config"
  end

  test do
    system "#{bin}/elitecode", "version"
  end
end