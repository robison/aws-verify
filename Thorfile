require 'thor'
require 'thor/actions'
require 'thor/scmversion'
require 'octokit'

## GitHub Helpers
module GitHub
  class << self
    def client
      raise 'Missing required environment variable GITHUB_TOKEN' unless env.include?('GITHUB_TOKEN')

      @client ||= Octokit::Client.new(:access_token => ENV['GITHUB_TOKEN'])
    end

    def verison
      IO.read(File.join(__dir__, 'VERSION'))
    end

    def repo
      env['TRAVIS_REPO_SLUG']
    end

    def commit
      ENV['TRAVIS_COMMIT']
    end
  end

  ## Thor Commands
  class Commands  < Thor
    namespace 'gh'

    desc 'tag', 'Use the GitHub API to create a tagged release for the current version'
    def tag
      say_staus :tag, "Creating draft release #{GitHub.version} on #{GitHub.repo}"
      GitHub.client.create_release(GitHub.repo, GitHub.version,
                                   :target_commitish => GitHub.commit,
                                   :draft => true)
    end

    desc 'upload BUILD_DIR=build', 'Upload all artifacts in BUILD_DIR to the current GitHub release'
    def upload(build_dir = 'build')

    end
  end
end

module Gox
  ## Thor Commands
  class Commands < Thor
    include Thor::Actions

    namespace 'gox'

    desc 'build BUILD_DIR=build', 'Perform a gox build, placing artifacts into BUILD_DIR'
    def build(build_dir = 'build')
      empty_directory build_dir
      run "gox -output=#{build_dir}/{{.Dir}}-{{.OS}}-{{.Arch}}"
    end
  end
end
