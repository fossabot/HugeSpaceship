package slots

import (
	"HugeSpaceship/pkg/common/db"
	"HugeSpaceship/pkg/common/model/auth"
	"HugeSpaceship/pkg/common/model/lbp_xml/slot"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func StartPublishHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		dbCtx := db.GetContext()
		s := slot.Slot{}
		err := ctx.BindXML(&s)
		if err != nil {
			log.Error().Err(err).Msg("Failed to parse xml body")
		}

		// This checks to see if the resources already exist in the DB
		c := 0
		resourcesToUpload := make([]string, 0, len(s.Resources))
		for i := range s.Resources {
			exists, err := db.ResourceExists(dbCtx, s.Resources[i])
			if err != nil {
				log.Warn().Err(err).Msg("failed to check if resource exists, assuming it doesn't")
			}
			if !exists {
				resourcesToUpload[c] = s.Resources[i]
				c++
			}
		}

		ctx.XML(200, slot.StartPublishSlotResponse{
			Resource: resourcesToUpload,
		})
	}
}

func PublishHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		dbCtx := db.GetContext()
		slotData := new(slot.Upload)
		err := ctx.BindXML(slotData)
		if err != nil {
			_ = ctx.Error(err)
			ctx.AbortWithStatus(400)
		}
		domain := ctx.GetInt("domain")
		session, _ := ctx.Get("session")

		id, err := db.InsertSlot(dbCtx, slotData, session.(auth.Session).UserID, domain)
		if err != nil {
			ctx.Error(err)
		}
		s, err := db.GetSlot(dbCtx, id)
		if err != nil {
			ctx.Error(err)
		}
		ctx.XML(200, &s)
	}
}
