# typed: false
# frozen_string_literal: true

# This file was generated by GoReleaser. DO NOT EDIT.
class Chainsaw < Formula
  desc "Declarative Kubernetes end-to-end testing."
  homepage "https://kyverno.github.io/chainsaw"
  version "0.2.8-beta.2"

  on_macos do
    on_intel do
      url "https://github.com/kyverno/chainsaw/releases/download/v0.2.8-beta.2/chainsaw_darwin_amd64.tar.gz"
      sha256 "2d7cc925df99f10bacaef7b223f89fa18cfdf6579409623a31cbae9500aa30b4"

      def install
        bin.install "chainsaw"
      end
    end
    on_arm do
      url "https://github.com/kyverno/chainsaw/releases/download/v0.2.8-beta.2/chainsaw_darwin_arm64.tar.gz"
      sha256 "130b6cb3bd312c783cd6306b310dca80af63c3dfbd5c5a99f96663145abd78f4"

      def install
        bin.install "chainsaw"
      end
    end
  end

  on_linux do
    on_intel do
      if Hardware::CPU.is_64_bit?
        url "https://github.com/kyverno/chainsaw/releases/download/v0.2.8-beta.2/chainsaw_linux_amd64.tar.gz"
        sha256 "9b9fef6205e6dda2ef6bb8a3044dce80f45c40b2df938380be2fa3971b17ed82"

        def install
          bin.install "chainsaw"
        end
      end
    end
    on_arm do
      if Hardware::CPU.is_64_bit?
        url "https://github.com/kyverno/chainsaw/releases/download/v0.2.8-beta.2/chainsaw_linux_arm64.tar.gz"
        sha256 "5ecb5f42c59798492e7ba3177dd43bedf48c967c800f55acf7f05ce996f41eb2"

        def install
          bin.install "chainsaw"
        end
      end
    end
  end
end
