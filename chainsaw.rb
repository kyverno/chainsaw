# typed: false
# frozen_string_literal: true

# This file was generated by GoReleaser. DO NOT EDIT.
class Chainsaw < Formula
  desc "Declarative Kubernetes end-to-end testing."
  homepage "https://kyverno.github.io/chainsaw"
  version "0.1.0"

  on_macos do
    if Hardware::CPU.intel?
      url "https://github.com/kyverno/chainsaw/releases/download/v0.1.0/chainsaw_darwin_amd64.tar.gz"
      sha256 "81abef255c3bd963956fc74e0bbc795dcb5ff55e9cdc11156deb870361004230"

      def install
        bin.install "chainsaw"
      end
    end
    if Hardware::CPU.arm?
      url "https://github.com/kyverno/chainsaw/releases/download/v0.1.0/chainsaw_darwin_arm64.tar.gz"
      sha256 "a9c36f34d6e468afa8a416d56819fdbfa3ee490ff47d589bc7fa5d934b33e968"

      def install
        bin.install "chainsaw"
      end
    end
  end

  on_linux do
    if Hardware::CPU.arm? && Hardware::CPU.is_64_bit?
      url "https://github.com/kyverno/chainsaw/releases/download/v0.1.0/chainsaw_linux_arm64.tar.gz"
      sha256 "5ecab357d70870a2a0b0f0b0aa9ab08dcd9349a808a88b8e96440f326b32d9ed"

      def install
        bin.install "chainsaw"
      end
    end
    if Hardware::CPU.intel?
      url "https://github.com/kyverno/chainsaw/releases/download/v0.1.0/chainsaw_linux_amd64.tar.gz"
      sha256 "620b84952c76e8498da5c8516a744f6bf89dbfede90530c39371331378d83d4d"

      def install
        bin.install "chainsaw"
      end
    end
  end
end
