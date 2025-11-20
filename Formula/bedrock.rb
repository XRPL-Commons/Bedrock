class Bedrock < Formula
  desc "XRPL smart contracts CLI tool"
  homepage "https://github.com/XRPL-Commons/Bedrock"
  url "https://github.com/XRPL-Commons/Bedrock/archive/refs/tags/v1.0.0.tar.gz"
  sha256 "df258a1d54201c650191ed0461aaffaf91caabda9ee79f1027bc14334ace4f2e"
  license "MIT"
  head "https://github.com/XRPL-Commons/Bedrock.git", branch: "main"

  depends_on "go" => :build

  def install
    system "go", "build", *std_go_args(ldflags: "-s -w -X main.version=#{version}"), "./cmd/bedrock"
  end

  test do
    assert_match "bedrock", shell_output("#{bin}/bedrock --help")
  end
end