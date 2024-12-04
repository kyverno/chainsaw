# typed: false
# frozen_string_literal: true

# This file was generated by GoReleaser. DO NOT EDIT.
class Chainsaw < Formula
  desc "Declarative Kubernetes end-to-end testing."
  homepage "https://kyverno.github.io/chainsaw"
  version "0.2.12"

  on_macos do
    on_intel do
      url "https://github.com/kyverno/chainsaw/releases/download/v0.2.12/chainsaw_darwin_amd64.tar.gz"
      sha256 "b49dba1214b32024567b1edc7653498a644fbef18111bcc3e1c46dc52e1d194e"

      def install
        bin.install "chainsaw"
      end
    end
    on_arm do
      url "https://github.com/kyverno/chainsaw/releases/download/v0.2.12/chainsaw_darwin_arm64.tar.gz"
      sha256 "717a07fcc4d781fff967b287880fed1d8b1e6af9fbecc7650a714c467f296d33"

      def install
        bin.install "chainsaw"
      end
    end
  end

  on_linux do
    on_intel do
      if Hardware::CPU.is_64_bit?
        url "https://github.com/kyverno/chainsaw/releases/download/v0.2.12/chainsaw_linux_amd64.tar.gz"
        sha256 "d6bfb17ba47af2db85edc0365288f92b5e1a4566f7ff130ec9b326f96856e209"

        def install
          bin.install "chainsaw"
        end
      end
    end
    on_arm do
      if Hardware::CPU.is_64_bit?
        url "https://github.com/kyverno/chainsaw/releases/download/v0.2.12/chainsaw_linux_arm64.tar.gz"
        sha256 "72a6273d6da16a04e29e0fae232631b084852d21ddf25f88ed3d3de480125d30"

        def install
          bin.install "chainsaw"
        end
      end
    end
  end
end
