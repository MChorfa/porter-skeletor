# Homebrew Formula for Skeletor - Porter Mixin Generator
class Skeletor < Formula
  desc "Enterprise-grade command-line tool for generating Porter mixins"
  homepage "https://getporter.github.io/skeletor/"
  url "https://github.com/getporter/skeletor/archive/v0.1.0.tar.gz"
  sha256 "PLACEHOLDER_SHA256"
  license "Apache-2.0"
  head "https://github.com/getporter/skeletor.git", branch: "main"

  depends_on "go" => :build

  def install
    # Set build variables
    ldflags = %W[
      -s -w
      -X github.com/getporter/skeletor/pkg/version.Version=#{version}
      -X github.com/getporter/skeletor/pkg/version.Commit=#{Utils.git_head}
      -X github.com/getporter/skeletor/pkg/version.Date=#{time.iso8601}
    ]

    # Build the binary
    system "go", "build", *std_go_args(ldflags: ldflags), "./cmd/skeletor"

    # Install shell completions
    generate_completions_from_executable(bin/"skeletor", "completion")

    # Install man page if available
    if File.exist?("docs/skeletor.1")
      man1.install "docs/skeletor.1"
    end
  end

  test do
    # Test version output
    assert_match version.to_s, shell_output("#{bin}/skeletor version")

    # Test help output
    help_output = shell_output("#{bin}/skeletor --help")
    assert_match "Create new Porter mixins easily", help_output

    # Test create command help
    create_help = shell_output("#{bin}/skeletor create --help")
    assert_match "Create a new Porter mixin", create_help
    assert_match "--enable-security", create_help
    assert_match "--enable-compliance", create_help
    assert_match "--enable-auth", create_help
    assert_match "--enable-observability", create_help

    # Test dry run (should not create files)
    system bin/"skeletor", "create", "--name", "test-mixin", "--author", "Test Author", "--dry-run", "--non-interactive"
    refute_predicate testpath/"test-mixin", :exist?

    # Test actual mixin creation
    system bin/"skeletor", "create", "--name", "test-mixin", "--author", "Test Author", "--non-interactive"
    assert_predicate testpath/"test-mixin", :exist?
    assert_predicate testpath/"test-mixin/go.mod", :exist?
    assert_predicate testpath/"test-mixin/README.md", :exist?
    assert_predicate testpath/"test-mixin/cmd/test-mixin", :exist?
    assert_predicate testpath/"test-mixin/pkg/test-mixin", :exist?

    # Test enterprise features
    system bin/"skeletor", "create", 
           "--name", "enterprise-mixin", 
           "--author", "Enterprise Author",
           "--enable-security",
           "--security-features", "input_validation,rate_limiting",
           "--enable-compliance",
           "--compliance-frameworks", "soc2",
           "--enable-auth",
           "--auth-features", "rbac",
           "--enable-observability",
           "--observability-features", "apm",
           "--non-interactive"
    
    assert_predicate testpath/"enterprise-mixin", :exist?
    assert_predicate testpath/"enterprise-mixin/pkg/security", :exist?
    assert_predicate testpath/"enterprise-mixin/pkg/compliance", :exist?
    assert_predicate testpath/"enterprise-mixin/pkg/auth", :exist?
    assert_predicate testpath/"enterprise-mixin/pkg/observability", :exist?
    assert_predicate testpath/"enterprise-mixin/configs/security.yaml", :exist?
    assert_predicate testpath/"enterprise-mixin/configs/compliance.yaml", :exist?
    assert_predicate testpath/"enterprise-mixin/configs/auth.yaml", :exist?
    assert_predicate testpath/"enterprise-mixin/configs/observability.yaml", :exist?
  end
end
