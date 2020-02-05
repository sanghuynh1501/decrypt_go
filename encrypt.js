let crypto;
try {
  crypto = require('crypto');
} catch (err) {
  console.log('crypto support is disabled!');
}

const ALGORITHM = 'aes-256-cbc';
const CIPHER_KEY = "abcdefghijklmnopqrstuvwxyz012345";  // Same key used in Golang
const BLOCK_SIZE = 16;

// Encrypts plain text into cipher text
function encrypt(plainText) {
  const iv = crypto.randomBytes(BLOCK_SIZE);
  const cipher = crypto.createCipheriv(ALGORITHM, CIPHER_KEY, iv);
  let cipherText;
  try {
    cipherText = cipher.update(plainText, 'utf8', 'hex');
    cipherText += cipher.final('hex');
    cipherText = iv.toString('hex') + cipherText
  } catch (e) {
    cipherText = null;
  }
  return cipherText;
}

const encrypt_text = encrypt("something something something")
console.log(encrypt_text)