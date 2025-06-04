class Rules < Formula
  desc "Internal rules CLI tool for Continue Dev"
  homepage "https://github.com/continuedev/rules-cli"
  version "0.0.1"
  
  # Since this is a private repo, we'll build from source
  
  depends_on "go" => :build

  def install
    # This assumes the user has Git access to the private repo
    system "git", "clone", "git@github.com:continuedev/rules-cli.git"
    cd "rules-cli" do
      # Optionally checkout a specific tag or branch
      system "git", "checkout", "main"
      
      # Build the binary
      system "go", "build", "-o", "rules", "-ldflags", "-X main.Version=#{version}", "./main.go"
      
      # Install the binary to the Homebrew prefix
      bin.install "rules"
    end
  end

  test do
    assert_match "#{version}", shell_output("#{bin}/rules --version")
  end
end