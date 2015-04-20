# -*- mode: python -*-
a = Analysis(['deis.py'],
             pathex=['.'],
             hiddenimports=[],
             hookspath=None,
             runtime_hooks=None)
pyz = PYZ(a.pure)
exe = EXE(pyz,
          a.scripts,
          a.binaries,
          a.zipfiles,
          a.datas,
          name='deis',
          debug=False,
          strip=None,
          upx=True,
          console=True)
