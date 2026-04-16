#!/usr/bin/env python3
"""
使用 Pillow 生成 H 字母风格的代码编辑器图标
"""
from pathlib import Path
from PIL import Image, ImageDraw

def create_icon(size, output_path):
    """创建指定尺寸的图标"""
    img = Image.new('RGBA', (size, size), (0, 0, 0, 0))
    draw = ImageDraw.Draw(img)
    
    # 绘制圆角矩形背景 - 蓝色
    padding = int(size * 0.03)
    corner_radius = int(size * 0.16)
    bg_color = (37, 99, 235)  # #2563eb
    
    draw.rounded_rectangle(
        [padding, padding, size-padding, size-padding],
        radius=corner_radius,
        fill=bg_color
    )
    
    # 绘制 H 字母 - 白色
    white = (255, 255, 255)
    stroke_width = max(3, size // 40)
    
    # 计算 H 的尺寸和位置
    h_width = int(size * 0.45)
    h_height = int(size * 0.5)
    h_left = (size - h_width) // 2
    h_top = (size - h_height) // 2 - int(size * 0.05)
    bar_width = max(12, h_width // 8)
    
    # 左竖线
    draw.rounded_rectangle(
        [h_left, h_top, h_left + bar_width, h_top + h_height],
        radius=bar_width // 2,
        fill=white
    )
    
    # 右竖线
    draw.rounded_rectangle(
        [h_left + h_width - bar_width, h_top, h_left + h_width, h_top + h_height],
        radius=bar_width // 2,
        fill=white
    )
    
    # 中间横线
    mid_y = h_top + (h_height - bar_width) // 2
    draw.rounded_rectangle(
        [h_left, mid_y, h_left + h_width, mid_y + bar_width],
        radius=bar_width // 2,
        fill=white
    )
    
    # 底部装饰线 - 浅蓝色
    accent_color = (147, 197, 253, 200)  # #93c5fd with alpha
    accent_y = int(size * 0.75)
    accent_width = int(size * 0.57)
    accent_height = max(2, size // 40)
    
    draw.rounded_rectangle(
        [(size - accent_width) // 2, accent_y,
         (size + accent_width) // 2, accent_y + accent_height],
        radius=accent_height // 2,
        fill=accent_color
    )
    
    # 保存
    img.save(output_path, 'PNG')
    return True

def generate_all_sizes(output_dir):
    """生成所有尺寸的图标"""
    output_dir = Path(output_dir)
    output_dir.mkdir(parents=True, exist_ok=True)
    
    # macOS 所需尺寸
    sizes = [16, 32, 64, 128, 256, 512, 1024]
    
    # 生成 iconset
    iconset_dir = output_dir / "AppIcon.iconset"
    if iconset_dir.exists():
        import shutil
        shutil.rmtree(iconset_dir)
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
    
    # 也生成一个标准的 appicon.png (1024x1024)
    appicon_path = output_dir.parent / "appicon.png"
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
    except FileNotFoundError:
        print("✗ iconutil 命令未找到，跳过 icns 生成")

if __name__ == "__main__":
    script_dir = Path(__file__).parent
    output_dir = script_dir / "darwin"
    
    print("使用 Pillow 生成 H 字母风格图标")
    print(f"输出目录: {output_dir}")
    print()
    
    generate_all_sizes(output_dir)
