#!/usr/bin/env python3
"""
生成不同尺寸的图标文件
"""
import sys
import subprocess
from pathlib import Path

def svg_to_png(svg_path, output_path, size):
    """使用 Inkscape 或 ImageMagick 将 SVG 转换为 PNG"""
    try:
        # 尝试使用 Inkscape
        subprocess.run([
            'inkscape', svg_path,
            f'--export-filename={output_path}',
            f'--export-width={size}',
            f'--export-height={size}'
        ], check=True, capture_output=True)
        return True
    except (subprocess.CalledProcessError, FileNotFoundError):
        pass
    
    try:
        # 尝试使用 ImageMagick
        subprocess.run([
            'magick', 'convert', '-density', '300',
            svg_path,
            '-resize', f'{size}x{size}',
            '-background', 'none',
            output_path
        ], check=True, capture_output=True)
        return True
    except (subprocess.CalledProcessError, FileNotFoundError):
        pass
    
    try:
        # 尝试使用 convert (旧版 ImageMagick)
        subprocess.run([
            'convert', '-density', '300',
            svg_path,
            '-resize', f'{size}x{size}',
            '-background', 'none',
            output_path
        ], check=True, capture_output=True)
        return True
    except (subprocess.CalledProcessError, FileNotFoundError):
        pass
    
    return False

def generate_iconset(svg_path, output_dir):
    """生成 macOS iconset"""
    iconset_dir = Path(output_dir) / "AppIcon.iconset"
    iconset_dir.mkdir(parents=True, exist_ok=True)
    
    sizes = [16, 32, 64, 128, 256, 512, 1024]
    
    for size in sizes:
        # 标准尺寸
        png_name = f"icon_{size}x{size}.png"
        png_path = iconset_dir / png_name
        if svg_to_png(svg_path, str(png_path), size):
            print(f"✓ 已生成 {png_name}")
        
        # 视网膜屏尺寸 (2x)
        if size <= 512:  # 1024 的 2x 太大，通常不需要
            png_2x_name = f"icon_{size}x{size}@2x.png"
            png_2x_path = iconset_dir / png_2x_name
            if svg_to_png(svg_path, str(png_2x_path), size * 2):
                print(f"✓ 已生成 {png_2x_name}")
    
    # 转换为 icns
    icns_path = Path(output_dir) / "icons.icns"
    try:
        subprocess.run([
            'iconutil', '-c', 'icns',
            str(iconset_dir),
            '-o', str(icns_path)
        ], check=True)
        print(f"✓ 已生成 icons.icns")
    except subprocess.CalledProcessError as e:
        print(f"✗ 生成 icns 失败: {e}")

if __name__ == "__main__":
    script_dir = Path(__file__).parent
    svg_path = sys.argv[1] if len(sys.argv) > 1 else str(script_dir / "editor-icon.svg")
    output_dir = sys.argv[2] if len(sys.argv) > 2 else str(script_dir / "darwin")
    
    print(f"从 SVG 生成图标: {svg_path}")
    print(f"输出目录: {output_dir}")
    print()
    
    generate_iconset(svg_path, output_dir)
