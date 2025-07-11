class Elitecode < Formula
  desc "CLI-focused coding platform for competitive programming"
  homepage "https://github.com/yourusername/elitecode"
  url "https://github.com/yourusername/elitecode/archive/refs/tags/v1.0.0.tar.gz"
  sha256 "your-sha256-hash-here"
  license "MIT"
  version "1.0.0"

  depends_on "node"
  depends_on "docker" => :optional

  def install
    # Install Node.js dependencies
    system "npm", "install", "--production"
    
    # Create libexec directory for the app
    libexec.install Dir["*"]
    
    # Create binary wrapper
    (bin/"elitecode").write_env_script("#{libexec}/bin/elitecode.js", {})
    
    # Install shell completions
    bash_completion.install "completions/elitecode.bash" => "elitecode"
    zsh_completion.install "completions/elitecode.zsh" => "_elitecode"
    fish_completion.install "completions/elitecode.fish"
  end

  def post_install
    # Create config directory
    config_dir = "#{Dir.home}/.elitecode"
    mkdir_p config_dir unless Dir.exist?(config_dir)
    
    # Set proper permissions
    File.chmod(0755, config_dir)
    
    ohai "EliteCode has been installed!"
    ohai "Run 'elitecode --help' to get started"
    ohai "Run 'elitecode init' to initialize your environment"
  end

  def caveats
    <<~EOS
      EliteCode CLI has been installed successfully!
      
      To get started:
        1. Run: elitecode init
        2. Run: elitecode signup (or elitecode login if you have an account)
        3. Run: elitecode problems to see available challenges
      
      Optional dependencies:
        - Docker: Required for running code in isolated containers
        - VS Code: Recommended for editing solutions
        - Git: Required for GitHub integration
      
      Configuration files are stored in: ~/.elitecode/
      
      For more information, visit: https://github.com/yourusername/elitecode
    EOS
  end

  test do
    # Test that the binary runs and shows help
    assert_match "EliteCode", shell_output("#{bin}/elitecode --help")
    
    # Test version command
    assert_match version.to_s, shell_output("#{bin}/elitecode --version")
    
    # Test that config directory can be created
    system "#{bin}/elitecode", "init", "--force"
    assert_predicate testpath/".elitecode", :exist?
  end
end