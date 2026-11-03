"""Build-time QR asset only. No wallet, RPC, payment or web QR service.

Development dependencies: qrcode, Pillow, pycryptodome and zxing-cpp.
The generated PNG contains exactly the owner-supplied address, not a payment
request. The network must be shown separately. No runtime dependency is added.
"""
import hashlib
import json
from pathlib import Path
import re

from Crypto.Hash import keccak
from PIL import Image
import qrcode
import zxingcpp

root = Path(__file__).resolve().parent.parent
config = json.loads((root / "release/sponsorship.json").read_text(encoding="utf-8"))
address = config["address"]
assert config["enabled"] and config["network"] == "Base" and config["chainId"] == 8453
assert re.fullmatch(r"0x[0-9a-fA-F]{40}", address)
hasher = keccak.new(digest_bits=256)
hasher.update(address[2:].lower().encode("ascii"))
digest = hasher.hexdigest()
checksum = "0x" + "".join(c.upper() if c in "abcdef" and int(digest[i], 16) >= 8 else c for i, c in enumerate(address[2:].lower()))
assert checksum == address, "Owner address checksum differs; do not silently change it"
assert config["qrPayload"] == "address-only" and config["asset"] is None and config["amount"] is None
qr = qrcode.QRCode(error_correction=qrcode.constants.ERROR_CORRECT_Q, box_size=10, border=4)
qr.add_data(address, optimize=0)
qr.make(fit=True)
target = root / "web/public/assets/support-base-address-qr.png"
qr.make_image(fill_color="black", back_color="white").convert("RGB").save(target)
decoded = zxingcpp.read_barcode(Image.open(target))
assert decoded and decoded.text == address, "Independent QR decode differs from configured address"
print(json.dumps({"checksumValid": True, "independentDecodeMatches": True, "network": "Base", "payload": "address-only", "sha256": hashlib.sha256(target.read_bytes()).hexdigest()}))
