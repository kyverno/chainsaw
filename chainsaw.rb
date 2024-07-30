# typed: false
# frozen_string_literal: true

# This file was generated by GoReleaser. DO NOT EDIT.
class Chainsaw < Formula
  desc "Declarative Kubernetes end-to-end testing."
  homepage "https://kyverno.github.io/chainsaw"
  version "0.2.8-beta.1"

  on_macos do
    on_intel do
      url "https://github.com/kyverno/chainsaw/releases/download/v0.2.8-beta.1/chainsaw_darwin_amd64.tar.gz"
      sha256 "5317f40e5d5603dd0123adb117d1ebf700be4f57a59c39e731e1a296d9935851"

      def install
        bin.install "chainsaw"
      end
    end
    on_arm do
      url "https://github.com/kyverno/chainsaw/releases/download/v0.2.8-beta.1/chainsaw_darwin_arm64.tar.gz"
      sha256 "18c87aacf74b4edfc58978aec66c110dccb2bbbf7182c002331b701cdbef2a44"

      def install
        bin.install "chainsaw"
      end
    end
  end

  on_linux do
    on_intel do
      if Hardware::CPU.is_64_bit?
        url "https://github.com/kyverno/chainsaw/releases/download/v0.2.8-beta.1/chainsaw_linux_amd64.tar.gz"
        sha256 "ab4f320ce3042f7f94a69e1e2881e6c05e4f740a30a035b8296277426e53cd72"

        def install
          bin.install "chainsaw"
        end
      end
    end
    on_arm do
      if Hardware::CPU.is_64_bit?
        url "https://github.com/kyverno/chainsaw/releases/download/v0.2.8-beta.1/chainsaw_linux_arm64.tar.gz"
        sha256 "647ca98fa22ffc95d3dc6a9147cbb6068b47ccf6410348cd7539abdac9b2fb37"

        def install
          bin.install "chainsaw"
        end
      end
    end
  end
end
