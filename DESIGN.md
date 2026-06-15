---
version: alpha
name: HackCulture
description: A bold, tech-forward innovation platform with bright electric accents, airy spacing, and clean editorial typography.
colors:
  primary: "#4953F5"
  primary-70: "#6A72F7"
  primary-40: "#AAB0FB"
  secondary: "#FFFFFF"
  tertiary: "#374151"
  neutral: "#F1F4F8"
  surface: "#FFFFFF"
  surface-muted: "#F7F8FC"
  on-surface: "#0F1729"
  border: "#E5E7EB99"
  overlay: "#121212"
  error: "#D92D20"
typography:
  headline-display:
    fontFamily: Space Grotesk
    fontSize: 56px
    fontWeight: 500
    lineHeight: 1.05
    letterSpacing: -0.03em
  headline-lg:
    fontFamily: Space Grotesk
    fontSize: 40px
    fontWeight: 500
    lineHeight: 43px
    letterSpacing: -0.02em
  headline-md:
    fontFamily: Space Grotesk
    fontSize: 33px
    fontWeight: 500
    lineHeight: 43px
  headline-sm:
    fontFamily: Space Grotesk
    fontSize: 28px
    fontWeight: 500
    lineHeight: 30px
  title-lg:
    fontFamily: Space Grotesk
    fontSize: 23px
    fontWeight: 500
    lineHeight: 28px
  body-lg:
    fontFamily: Space Grotesk
    fontSize: 19px
    fontWeight: 400
    lineHeight: 31px
  body-md:
    fontFamily: Space Grotesk
    fontSize: 16px
    fontWeight: 400
    lineHeight: 26px
  body-sm:
    fontFamily: Space Grotesk
    fontSize: 14px
    fontWeight: 400
    lineHeight: 22px
  label-lg:
    fontFamily: Space Grotesk
    fontSize: 18px
    fontWeight: 700
    lineHeight: 1.2
  label-md:
    fontFamily: Space Grotesk
    fontSize: 16px
    fontWeight: 700
    lineHeight: 1.2
  label-sm:
    fontFamily: Space Grotesk
    fontSize: 12px
    fontWeight: 600
    lineHeight: 1.2
    letterSpacing: 0.08em
rounded:
  none: 0px
  sm: 4px
  md: 8px
  lg: 16px
  xl: 24px
  full: 9999px
spacing:
  xs: 6px
  sm: 16px
  md: 28px
  lg: 40px
  xl: 58px
components:
  button-primary:
    backgroundColor: "{colors.secondary}"
    textColor: "{colors.tertiary}"
    typography: "{typography.label-lg}"
    rounded: "{rounded.md}"
    padding: 8px 16px
    width: 140px
    height: 44px
  button-primary-hover:
    backgroundColor: "{colors.surface-muted}"
    textColor: "{colors.tertiary}"
    typography: "{typography.label-lg}"
    rounded: "{rounded.md}"
  button-secondary:
    backgroundColor: "transparent"
    textColor: "{colors.secondary}"
    typography: "{typography.label-lg}"
    rounded: "{rounded.sm}"
    padding: 8px 16px
    width: 140px
    height: 44px
  button-link:
    backgroundColor: "transparent"
    textColor: "{colors.secondary}"
    typography: "{typography.body-md}"
    rounded: "{rounded.none}"
  card:
    backgroundColor: "{colors.neutral}"
    textColor: "{colors.on-surface}"
    typography: "{typography.body-md}"
    rounded: "{rounded.lg}"
    padding: 28px
  input:
    backgroundColor: "{colors.surface}"
    textColor: "{colors.on-surface}"
    typography: "{typography.body-md}"
    rounded: "{rounded.md}"
    padding: 12px 16px
  chip:
    backgroundColor: "{colors.primary-40}"
    textColor: "{colors.secondary}"
    typography: "{typography.label-sm}"
    rounded: "{rounded.full}"
    padding: 6px 12px
---

# HackCulture

## Overview
HackCulture feels like a polished innovation-stage brand: energetic, optimistic, and designed to inspire corporate teams without feeling stiff. The visual tone is modern and slightly futuristic, with a bright electric gradient that shifts from blue into violet and a very clean white content area below. It balances spacious hero composition with structured, enterprise-friendly navigation and calls to action, making it suitable for innovation programs, hackathons, and capability-building audiences.

## Colors
- **Primary (#4953F5):** The electric blue-violet brand accent used for the hero gradient, emphasis, and interactive highlights. It gives the page its energetic, startup-adjacent personality.
- **Secondary (#FFFFFF):** The main text and button fill color on dark and gradient surfaces. It keeps the interface crisp and high-contrast.
- **Tertiary (#374151):** A subdued dark neutral used for readable button text and supporting UI on light surfaces. It feels professional rather than decorative.
- **Neutral (#F1F4F8):** A soft cool gray used for cards and section surfaces. It keeps content blocks light while avoiding stark white everywhere.
- **Surface (#FFFFFF):** The pure base surface for the lower page area and primary button backgrounds.
- **Surface-muted (#F7F8FC):** A near-white background for subtle layering and secondary containers.
- **On-surface (#0F1729):** The deep ink tone for body content and card text, providing strong readability.
- **Border (#E5E7EB99):** A faint translucent border used to define cards and modules without introducing heavy outlines.
- **Overlay (#121212):** The dark mode background reference from the extracted styleguide; useful for deep sections or modal backdrops.
- **Error (#D92D20):** Reserved for validation and destructive states; not prominent in the screenshot, but needed for system completeness.

## Typography
The system uses Space Grotesk throughout, which keeps the voice contemporary, geometric, and slightly editorial. Headings are medium weight, with the largest hero text using strong negative letter spacing for a compact, cinematic feel. Body text stays lighter and more open, while labels and buttons use bold weights to make actions stand out clearly. Uppercase treatments appear in the hero and section headings, reinforcing the structured, campaign-like presentation.

- **Headlines:** Use Space Grotesk at medium weight for strong hierarchy. Hero-scale display text should be large, tight, and visually centered.
- **Body:** Use Space Grotesk regular for descriptive copy and supporting statements. Line height should remain generous to preserve the airy feel.
- **Labels:** Use bold Space Grotesk for buttons, navigation, and emphasis. Small labels may use added letter spacing for utility-style clarity.

## Layout
The page uses a wide, centered layout with a full-bleed hero that stretches edge to edge and a more conventional content rhythm below. Spacing feels expansive: large vertical gaps in the hero, then structured section padding and evenly distributed cards. The extracted spacing scale of 6px, 16px, 28px, 40px, and 58px suggests a modest modular rhythm, with 16px as the base unit and larger jumps for section separation. Cards and feature blocks rely on consistent internal padding rather than dense nested layouts.

## Elevation & Depth
Depth is subtle and controlled. The design relies more on color layering, faint borders, and soft shadows than on dramatic elevation. The hero gradient creates the primary sense of depth, while cards use a light gray fill, a delicate translucent border, and a soft shadow to separate them from the page without feeling heavy. Overall, the system is more tonal than shadow-driven.

## Shapes
The shape language is friendly and modern, with medium radii on buttons and cards. Buttons use 8px and 4px rounding depending on emphasis, while cards sit at a larger 16px radius for softer modular panels. Full-pill chips and tags should use the `full` radius to preserve the rounded, approachable aesthetic. Avoid overly sharp corners on interactive elements unless a utility link is intended.

## Components
- **Primary buttons (`button-primary`):** White-filled CTA buttons with dark text, bold label styling, 140px minimum width, and 44px height. They should feel prominent and inviting, with medium rounding and compact horizontal padding.
- **Secondary buttons (`button-secondary`):** Transparent buttons with white outlines and white text, used on the gradient hero. They are visually lighter than the primary CTA but still clearly interactive.
- **Hover states (`button-primary-hover`):** Keep hover changes subtle, using a slight tonal shift rather than strong shadows or movement.
- **Link buttons (`button-link`):** Simple text links with no container, no border, and underlined styling. Use them for low-emphasis actions in navigation or footers.
- **Cards (`card`):** Large, light-toned panels with 28px padding, 16px rounding, faint borders, and a soft shadow. They should support image-first or content-first layouts without feeling boxed in.
- **Inputs (`input`):** Controls should mirror the card language: clean surface fill, 8px rounding, ample padding, and clear legibility. Borders should remain subtle rather than dark or dense.
- **Chips (`chip`):** Small pill shapes for status or category tags, using the `full` radius and compact padding. They should read as secondary metadata, not as primary actions.

## Do's and Don'ts
- Do keep hero sections spacious, centered, and high-contrast.
- Do use Space Grotesk consistently across headings, UI labels, and body copy.
- Do favor white, light gray, and the electric blue-violet accent as the core palette.
- Do preserve subtle depth with faint borders and restrained shadows.
- Don't introduce heavy drop shadows or glossy effects.
- Don't use dense paragraph blocks or cramped spacing.
- Don't replace the clean geometric type with decorative or highly stylized fonts.
- Don't make secondary actions visually louder than the primary CTA.