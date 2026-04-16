#!/usr/bin/env python3
"""
使用 Pillow 生成图标
安装: pip install Pillow
"""
import sys
from pathlib import Path
from PIL import Image, ImageDraw, ImageFont

def create_icon(size, output_path):
    """创建指定尺寸的图标"""
    # 创建图像
    img = Image.new('RGBA', (size, size), (0, 0, 0, 0))
    draw = ImageDraw.Draw(img)
    
    # 绘制圆角矩形背景
    padding = int(size * 0.03)
    corner_radius = int(size * 0.16)
    
    # 渐变背景 - 简化为纯色
    bg_color = (37, 99, 235)  # #2563eb
    draw.rounded_rectangle(
        [padding, padding, size-padding, size-padding],
        radius=corner_radius,
        fill=bg_color
    )
    
    # 绘制代码符号 < />
    # 计算位置
    center_x = size // 2
    center_y = size // 2
    symbol_width = int(size * 0.6)
    symbol_height = int(size * 0.5)
    
    # 左侧 <
    left_x = center_x - symbol_width // 3
    left_top = center_y - symbol_height // 2
    left_bottom = center_y + symbol_height // 2
    
    draw.line(
        [(left_x, left_top), 
         (left_x + symbol_width // 6, center_y),
         (left_x, left_bottom)],
        fill='white',
        width=max(3, size // 20)
    )
    
    # 中间 /
    mid_x = center_x
    draw.line(
        [(mid_x - symbol_width // 12, left_top),
         (mid_x + symbol_width // 12, left_bottom)],
        fill='white',
        width=max(3, size // 20)
    )
    
    # 右侧 >
    right_x = center_x + symbol_width // 3
    draw.line(
        [(right_x, left_top),
         (right_x - symbol_width // 6, center_y),
         (right_x, left_bottom)],
        fill='white',
        width=max(3, size // 20)
    )
    
    # 底部装饰线
    line_y = int(size * 0.75)
    line_width = int(size * 0.6)
    line_height = max(2, size // 40)
    draw.rounded_rectangle(
        [center_x - line_width // 2, line_y,
         center_x + line_width // 2, line_y + line_height],
        radius=line_height // 2,
        fill=(147, 197, 253, 200)  # 半透明浅蓝色
    )
    
    # 保存
    img.save(output_path, 'PNG')
    return True

def generate_all_sizes(svg_path, output_dir):
    """生成所有尺寸的图标"""
    output_dir = Path(output_dir)
    output_dir.mkdir(parents=True, exist_ok=True)
    
    # macOS 所需尺寸
    sizes = [16, 32, 64, 128, 256, 512, 1024]
    
    # 生成 iconset
    iconset_dir = output_dir / "AppIcon.iconset"
    iconset_dir.mkdir(parents=True, exist_ok=True)
    
    for size in sizes:
        # 标准尺寸
        png_path = iconset_dir / f"icon_{size}x{size}.png"
        if create_icon(size, str(png_path)):
            print(f"✓ 已生成 icon_{size}x{size}.png")
        
        # 2x 尺寸 (除了 1024)
        if size <= 512:
            png_2x_path = iconset_dir / f"icon_{size}x{size}@2x.png"
            if create_icon(size * 2, str(png_2x_path)):
                print(f"✓ 已生成 icon_{size}x{size}@2x.png")
    
    # 也生成一个标准的 appicon.png
    appicon_path = output_dir / "appicon.png"
    if create_icon(1024, str(appicon_path)):
        print(f"✓ 已生成 appicon.png")
    
    # 转换为 icns (macOS)
    import subprocess
    try:
        subprocess.run([
            'iconutil', '-c', 'icns',
            str(iconset_dir),
            '-o', str(output_dir / "icons.icns")
        ], check=True)
        print(f"✓ 已生成 icons.icns")
    except subprocess.CalledProcessError as e:
        print(f"✗ 生成 icns 失败: {e}")
        print("  提示: 在 macOS 上运行此脚本以生成 .icns 文件")

if __name__ == "__main__":
    script_dir = Path(__file__).parent
    output_dir = script_dir / "darwin"
    
    print("使用 Pillow 生成图标")
    print(f"输出目录: {output_dir}")
    print()
    
    generate_all_sizes(None, output_dir)
