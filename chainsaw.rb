# typed: false
# frozen_string_literal: true

# This file was generated by GoReleaser. DO NOT EDIT.
class Chainsaw < Formula
  desc "Declarative Kubernetes end-to-end testing."
  homepage "https://kyverno.github.io/chainsaw"
  version "0.1.9"

  on_macos do
    if Hardware::CPU.arm?
      url "https://github.com/kyverno/chainsaw/releases/download/v0.1.9/chainsaw_darwin_arm64.tar.gz"
      sha256 "4187956ba26fd5dadf6552bfa77e769afdac48c08e5741c46a8e38b07ca708bc"

      def install
        bin.install "chainsaw"
      end
    end
    if Hardware::CPU.intel?
      url "https://github.com/kyverno/chainsaw/releases/download/v0.1.9/chainsaw_darwin_amd64.tar.gz"
      sha256 "8c65f3ee952aa04754d644f2ef3d5f489153638de4e71de6348d4628e5af0378"

      def install
        bin.install "chainsaw"
      end
    end
  end

  on_linux do
    if Hardware::CPU.intel?
      url "https://github.com/kyverno/chainsaw/releases/download/v0.1.9/chainsaw_linux_amd64.tar.gz"
      sha256 "4080d3bb5ea6de6f85198e413e24a5c7aee941f027ba8b545f7a1ddbaa2e2856"

      def install
        bin.install "chainsaw"
      end
    end
    if Hardware::CPU.arm? && Hardware::CPU.is_64_bit?
      url "https://github.com/kyverno/chainsaw/releases/download/v0.1.9/chainsaw_linux_arm64.tar.gz"
      sha256 "05f2cdce3f34989e71f47cd30e005a49fa8d7abefbede20311f96eed016a34b8"

      def install
        bin.install "chainsaw"
      end
    end
  end
end
