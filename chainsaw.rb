# typed: false
# frozen_string_literal: true

# This file was generated by GoReleaser. DO NOT EDIT.
class Chainsaw < Formula
  desc "Declarative Kubernetes end-to-end testing."
  homepage "https://kyverno.github.io/chainsaw"
  version "0.1.8"

  on_macos do
    if Hardware::CPU.arm?
      url "https://github.com/kyverno/chainsaw/releases/download/v0.1.8/chainsaw_darwin_arm64.tar.gz"
      sha256 "41f1d7992ef266a95273e371e98c2bc835c0c310d2b299b0f44852c64ad0569c"

      def install
        bin.install "chainsaw"
      end
    end
    if Hardware::CPU.intel?
      url "https://github.com/kyverno/chainsaw/releases/download/v0.1.8/chainsaw_darwin_amd64.tar.gz"
      sha256 "df46ae522dcb6e9f9d13530e19b808a9d3f293dc22018b918ab26dc40545f184"

      def install
        bin.install "chainsaw"
      end
    end
  end

  on_linux do
    if Hardware::CPU.arm? && Hardware::CPU.is_64_bit?
      url "https://github.com/kyverno/chainsaw/releases/download/v0.1.8/chainsaw_linux_arm64.tar.gz"
      sha256 "476af662222e4842f08fb21055b2260a5fed3c2e3ccc1ce3fd6eba51cae64254"

      def install
        bin.install "chainsaw"
      end
    end
    if Hardware::CPU.intel?
      url "https://github.com/kyverno/chainsaw/releases/download/v0.1.8/chainsaw_linux_amd64.tar.gz"
      sha256 "f720cc50dc53b6cedf17a9630dbc12737d5ac69a9dccd6623b5d8ea2f2b210fc"

      def install
        bin.install "chainsaw"
      end
    end
  end
end
