class ElitecodeBrew < Formula
  desc "A coding practice platform"
  homepage "https://github.com/IshaanNene/EliteCode-brew"
  url "https://github.com/IshaanNene/EliteCode-brew/archive/refs/tags/v0.1.0.tar.gz"
  sha256 "" # Replace with the actual SHA-256 hash of the tarball
  license "MIT"

  depends_on "go" => :build

  def install
    system "go", "build", "-o", "elitecode-brew"
    bin.install "elitecode-brew"
  end

  test do
    system "#{bin}/elitecode-brew", "--version"
  end
end 