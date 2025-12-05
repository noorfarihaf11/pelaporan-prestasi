package repository

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetAchievementStatistics(db *mongo.Database) (map[string]interface{}, error) {
	collection := db.Collection("achievements")

	ctx := context.Background()

	// Total prestasi per tipe
	typeCountCursor, err := collection.Aggregate(ctx, bson.A{
		bson.M{"$group": bson.M{"_id": "$achievement_type", "count": bson.M{"$sum": 1}}},
	})
	if err != nil {
		return nil, err
	}
	var typeCounts []bson.M
	if err := typeCountCursor.All(ctx, &typeCounts); err != nil {
		return nil, err
	}

	// Total prestasi per periode (tahun dari created_at)
	periodCountCursor, err := collection.Aggregate(ctx, bson.A{
		bson.M{"$group": bson.M{"_id": bson.M{"year": bson.M{"$year": "$created_at"}}, "count": bson.M{"$sum": 1}}},
	})
	if err != nil {
		return nil, err
	}
	var periodCounts []bson.M
	if err := periodCountCursor.All(ctx, &periodCounts); err != nil {
		return nil, err
	}

	// Top mahasiswa berprestasi (sum points)
	topStudentsCursor, err := collection.Aggregate(ctx, bson.A{
		bson.M{"$group": bson.M{"_id": "$student_id", "total_points": bson.M{"$sum": "$points"}}},
		bson.M{"$sort": bson.M{"total_points": -1}},
		bson.M{"$limit": 10},
	})
	if err != nil {
		return nil, err
	}
	var topStudents []bson.M
	if err := topStudentsCursor.All(ctx, &topStudents); err != nil {
		return nil, err
	}

	// Distribusi tingkat kompetisi
	levelCursor, err := collection.Aggregate(ctx, bson.A{
		bson.M{"$group": bson.M{"_id": "$details.level", "count": bson.M{"$sum": 1}}},
	})
	if err != nil {
		return nil, err
	}
	var levelDistribution []bson.M
	if err := levelCursor.All(ctx, &levelDistribution); err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"type_counts":       typeCounts,
		"period_counts":     periodCounts,
		"top_students":      topStudents,
		"level_distribution": levelDistribution,
	}, nil
}
func GetStudentStatistics(db *mongo.Database, studentID string) (map[string]interface{}, error) {
	collection := db.Collection("achievements")
	ctx := context.Background()

	matchStage := bson.M{"$match": bson.M{"student_id": studentID}}

	// Total prestasi per tipe
	typeCountCursor, err := collection.Aggregate(ctx, bson.A{
		matchStage,
		bson.M{"$group": bson.M{"_id": "$achievement_type", "count": bson.M{"$sum": 1}}},
	})
	if err != nil {
		return nil, err
	}
	var typeCounts []bson.M
	if err := typeCountCursor.All(ctx, &typeCounts); err != nil {
		return nil, err
	}

	// Total prestasi per periode (tahun)
	periodCountCursor, err := collection.Aggregate(ctx, bson.A{
		matchStage,
		bson.M{"$group": bson.M{"_id": bson.M{"year": bson.M{"$year": "$created_at"}}, "count": bson.M{"$sum": 1}}},
	})
	if err != nil {
		return nil, err
	}
	var periodCounts []bson.M
	if err := periodCountCursor.All(ctx, &periodCounts); err != nil {
		return nil, err
	}

	// Total points
	pointsCursor, err := collection.Aggregate(ctx, bson.A{
		matchStage,
		bson.M{"$group": bson.M{"_id": "$student_id", "total_points": bson.M{"$sum": "$points"}}},
	})
	if err != nil {
		return nil, err
	}
	var totalPoints []bson.M
	if err := pointsCursor.All(ctx, &totalPoints); err != nil {
		return nil, err
	}

	// Distribusi tingkat kompetisi
	levelCursor, err := collection.Aggregate(ctx, bson.A{
		matchStage,
		bson.M{"$group": bson.M{"_id": "$details.level", "count": bson.M{"$sum": 1}}},
	})
	if err != nil {
		return nil, err
	}
	var levelDistribution []bson.M
	if err := levelCursor.All(ctx, &levelDistribution); err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"type_counts":       typeCounts,
		"period_counts":     periodCounts,
		"total_points":      totalPoints,
		"level_distribution": levelDistribution,
	}, nil
}
