package shopping_guide

import "github.com/dcsunny/wechat/context"

type Guide struct {
	*context.Context
	GuideManager  *GuideManager
	GuideBuyer    *GuideBuyer
	GuideTag      *GuideTag
	GuideMass     *GuideMass
	GuideMaterial *GuideMaterial
}

func NewGuide(ctx *context.Context) *Guide {
	guide := new(Guide)
	guide.GuideManager = NewGuideManager(ctx)
	guide.GuideBuyer = NewGuideBuyer(ctx)
	guide.GuideTag = NewGuideTag(ctx)
	guide.GuideMass = NewGuideMass(ctx)
	guide.GuideMaterial = NewGuideMaterial(ctx)
	return guide
}
