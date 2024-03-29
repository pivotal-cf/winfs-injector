# typed: false
# frozen_string_literal: true

# This file was generated by GoReleaser. DO NOT EDIT.
class WinfsInjector < Formula
  desc ""
  homepage ""
  version "0.25.0"

  on_macos do
    url "https://github.com/pivotal-cf/winfs-injector/releases/download/0.25.0/winfs-injector-darwin.tar.gz"
    sha256 "f1e3920dd0429fdbcce7a68a68766de3ff21408be614ab3c0c3d8a631092f9e4"

    def install
      bin.install "winfs-injector"
    end

    if Hardware::CPU.arm?
      def caveats
        <<~EOS
          The darwin_arm64 architecture is not supported for the WinfsInjector
          formula at this time. The darwin_amd64 binary may work in compatibility
          mode, but it might not be fully supported.
        EOS
      end
    end
  end

  on_linux do
    if Hardware::CPU.intel?
      url "https://github.com/pivotal-cf/winfs-injector/releases/download/0.25.0/winfs-injector-linux.tar.gz"
      sha256 "a8dfde2cb4f67998b01d41468d26ae9d823fe3a635d8cbd4601aa00e10e1b535"

      def install
        bin.install "winfs-injector"
      end
    end
  end

  test do
    system "#{bin}/winfs-injector --version"
  end
end
