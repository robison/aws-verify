'use strict';
const Assert = require('assert');
const Crypto = require('crypto');
const CP = require('child_process');
const OS = require('os');
const Path = require('path');

const Logger = require('2x4').export(module);

/**
 * Manage an instance of the aws-verify server
 */
class Runner {
  /**
   * @constructor
   */
  constructor(options) {
    super();

    options = Object.assign({
      maxRetries: 5,
      backoff: 0.5,
      certificates: []
    }, options);

    this.arch = Runner.ARCH[process.arch];
    this.platform = Runner.PLATFORM[process.platform];
    this.certificates = options.certificates;

    Assert.ok(this.arch, `Unsupported architecture ${process.arch}`);
    Assert.ok(this.platform, `Unsupported platform ${process.platform}`);
    Assert.ok(this.certificates instanceof Array, 'Parameter `certificates` must be an Array');

    const random = Crypto.randomBytes(16).toString('hex');

    this.socket = Path.join(OS.tmpdir(), `${process.pid}-${random}.sock`);
    this.logger = Logger.child({});

    this.stopped = false;
    this.maxRetries = options.maxRetries;
    this.backoff = options.backoff;
    this.retires = 0;
  }

  /**
   * Get the path to the binary for the current version/platform
   * @return {String}
   */
  get binary() {
    return Path.join(Runner.BIN,
      `aws-verify-${Runner.VERSION}-${this.platform}-${this.arch}`);
  }

  /**
   * Start the service
   */
  start() {
    if (this.process) { return; }
    this.stopped = false;

    const args = [`-socket=${this.socket}`];

    if (this.certificates.length > 0) {
      args.push(`-certificates=${this.certificates.join(',')}`);
    }

    this.process = CP.spawn(this.binary, args);
    this.process.on('close', (code) => {
      delete this.process;
      if (this.stopped) { return; }
    });

    this.process.stdout.on('data', (message) => )
  }

  /**
   * Stop the service
   */
  stop() {
    this.stopped = true;
    if (!this.process) { return; }

    this.process.kill('SIGTERM');
  }
}

/**
 * Absolute path to the package's binary directory
 * @type {String}
 */
Runner.BIN = Path.resolve(__dirname, '../bin');

/**
 * The package version
 * @type {String}
 */
Runner.VERSION = require('../package').version;

/**
 * Mappings from supported values of `process.arch` to the
 * corresponding Golang arch names
 * @type {Object}
 */
Runner.ARCH = {
  arm: 'arm',
  x64: 'amd64',
  i386: '386'
};

/**
 * Mappings from supported values of `process.platform` to the
 * corresponding Golang platform names
 * @type {Object}
 */
Runner.PLATFORM = {
  darwin: 'darwin',
  linux: 'linux'
};
